package db

import (
	"context"
	"fmt"
	"log"
	"projet/internal/config"
	"projet/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	//10seconde bech tsir connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//baad ma tufa fnct taaml cancel 9bl 10sec
	defer cancel()

	//tkhou les informations il config
	dsn := cfg.GetDBConnectionString()

	//yssir affichage fi terminal ll requests sql
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	//houni tssir l connexion wl cas d'erreur
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	//l Auto migration ti creati table si n'xiste pas. Automatiquement.
	db.AutoMigrate(&models.Produit{})

	//récupération de la connexion behind the scenes.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %v", err)
	}

	//la bdd ou réponds ou pas.
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	//succès
	log.Println("Connected to PostgreSQL successfully!")
	return db, nil
}

func Disconnect(db *gorm.DB) error {
	//idhekeni nil ma3andik matskkr
	if db == nil {
		return nil
	}
	// get raw sql connection that gorm is using.
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %v", err)
	}
	//tskkr
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to disconnect from PostgreSQL: %v", err)
	}
	//success
	log.Println(" Disconnected from PostgreSQL")
	return nil
}
