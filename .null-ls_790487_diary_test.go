package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Test_addDiaryEntry(t *testing.T) {
	type args struct {
		db *gorm.DB
		c  *gin.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addDiaryEntry(tt.args.db, tt.args.c)
		})
	}
}

func Test_getDiaryEntry(t *testing.T) {
	type args struct {
		db *gorm.DB
		c  *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getDiaryEntry(tt.args.db, tt.args.c)
		})
	}
}

func Test_updateDiaryEntry(t *testing.T) {
	type args struct {
		db *gorm.DB
		c  *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateDiaryEntry(tt.args.db, tt.args.c)
		})
	}
}

func Test_deleteDiaryEntry(t *testing.T) {
	type args struct {
		db *gorm.DB
		c  *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteDiaryEntry(tt.args.db, tt.args.c)
		})
	}
}
