package main

import (
	"flag"
	"log"

	"go-music-streamer/internal/config"
	"go-music-streamer/internal/database"
	"go-music-streamer/internal/database/postgres"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	var targetUserID uint
	flag.UintVar(&targetUserID, "userid", 1, "Target User ID to associate playlists and songs")
	flag.Parse()

	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Connecting to databases...")
	dbClient, err := database.ConnectDBSources(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to databases: %v", err)
	}

	pgDB, ok := dbClient.GetPostgres().(*gorm.DB)
	if !ok || pgDB == nil {
		log.Fatalf("Failed to retrieve Postgres connection")
	}

	log.Println("Migrating Database schemas...")
	err = pgDB.AutoMigrate(
		&postgres.User{},
		&postgres.UserRole{},
		&postgres.Role{},
		&postgres.RolePermission{},
		&postgres.Resource{},
		&postgres.Action{},
		&postgres.Admin{},
		&postgres.Playlist{},
		&postgres.Song{},
		&postgres.UserLikedSong{},
		&postgres.AdminRole{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schemas: %v", err)
	}

	log.Println("Starting database seeding...")

	// 1. Initializing Admin
	log.Println("Seeding Admin...")
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("adminpass123"), bcrypt.DefaultCost)
	admin := postgres.Admin{
		Email:    "admin@musicstreamer.com",
		Password: string(adminHash),
	}
	pgDB.Where(postgres.Admin{Email: admin.Email}).FirstOrCreate(&admin)

	// 2. Initializing Actions
	log.Println("Seeding Actions...")
	actions := []postgres.Action{
		{Name: "CREATE"}, {Name: "READ"}, {Name: "UPDATE"}, {Name: "DELETE"},
	}
	for i, action := range actions {
		pgDB.Where(postgres.Action{Name: action.Name}).FirstOrCreate(&actions[i])
	}

	// 3. Initializing Resources
	log.Println("Seeding Resources...")
	resources := []postgres.Resource{
		{Name: "SONG"}, {Name: "PLAYLIST"}, {Name: "USER"}, {Name: "ROLE"},
	}
	for i, res := range resources {
		pgDB.Where(postgres.Resource{Name: res.Name}).FirstOrCreate(&resources[i])
	}

	// 4. Initializing Roles
	log.Println("Seeding Roles...")
	roles := []postgres.Role{
		{Name: "SUPER_ADMIN"}, {Name: "CREATOR"}, {Name: "LISTENER"},
	}
	for i, role := range roles {
		pgDB.Where(postgres.Role{Name: role.Name}).FirstOrCreate(&roles[i])
	}

	// 5. Initializing Permissions (Assign all to SUPER_ADMIN for demonstration)
	log.Println("Seeding Permissions & RolePermissions...")
	for _, res := range resources {
		for _, act := range actions {
			perm := postgres.RolePermission{
				RoleID:     roles[0].ID,
				ResourceID: res.ID,
				ActionID:   act.ID,
			}
			pgDB.Where(postgres.RolePermission{
				RoleID:     roles[0].ID,
				ResourceID: res.ID,
				ActionID:   act.ID,
			}).FirstOrCreate(&perm)
		}
	}

	// 6. Ensure configuring User exists
	log.Println("Checking for target User ID", targetUserID, 10)
	targetUser := postgres.User{}
	res := pgDB.First(&targetUser, targetUserID)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			userHash, _ := bcrypt.GenerateFromPassword([]byte("userpass123"), bcrypt.DefaultCost)
			targetUser = postgres.User{
				Username:  "seeduser",
				Email:     "seeduser@musicstreamer.com",
				Password:  string(userHash),
				IsCreator: true,
			}
			targetUser.ID = targetUserID
			pgDB.Where(postgres.User{Email: targetUser.Email}).FirstOrCreate(&targetUser)

			// Try attaching role
			roleMapping := postgres.UserRole{
				UserID: targetUser.ID,
				RoleID: roles[1].ID, // CREATOR
			}
			pgDB.Where(postgres.UserRole{UserID: targetUser.ID, RoleID: roles[1].ID}).FirstOrCreate(&roleMapping)
		} else {
			log.Fatalf("Error querying user: %v", res.Error)
		}
	}

	// 7. Seed Songs
	log.Println("Seeding Songs...")
	songs := []postgres.Song{
		{Title: "Bohemian Rhapsody", Artist: "Queen", Album: "A Night at the Opera", Duration: 354, URL: "http://example.com/song1.mp3"},
		{Title: "Stairway to Heaven", Artist: "Led Zeppelin", Album: "Led Zeppelin IV", Duration: 482, URL: "http://example.com/song2.mp3"},
		{Title: "Hotel California", Artist: "Eagles", Album: "Hotel California", Duration: 390, URL: "http://example.com/song3.mp3"},
		{Title: "Sweet Child O' Mine", Artist: "Guns N' Roses", Album: "Appetite for Destruction", Duration: 356, URL: "http://example.com/song4.mp3"},
		{Title: "Imagine", Artist: "John Lennon", Album: "Imagine", Duration: 183, URL: "http://example.com/song5.mp3"},
		{Title: "Comfortably Numb", Artist: "Pink Floyd", Album: "The Wall", Duration: 382, URL: "http://example.com/song6.mp3"},
		{Title: "Hey Jude", Artist: "The Beatles", Album: "Hey Jude", Duration: 431, URL: "http://example.com/song7.mp3"},
		{Title: "Smells Like Teen Spirit", Artist: "Nirvana", Album: "Nevermind", Duration: 301, URL: "http://example.com/song8.mp3"},
	}
	for i, song := range songs {
		pgDB.Where(postgres.Song{Title: song.Title, Artist: song.Artist}).FirstOrCreate(&songs[i])
	}

	// 8. Seed Playlist & Associate Songs
	log.Printf("Seeding Playlists for User %d...\n", targetUser.ID)
	playlist := postgres.Playlist{
		Name:       "Rock Classics Seeded",
		Visibility: "public",
		CreatedBy:  targetUser.ID,
	}
	pgDB.Where(postgres.Playlist{Name: playlist.Name, CreatedBy: targetUser.ID}).FirstOrCreate(&playlist)

	log.Println("Associating Songs to Playlist...")
	for _, song := range songs {
		pgDB.Model(&playlist).Association("Songs").Append(&song)
	}

	log.Println("Database seeding completed successfully.")
}
