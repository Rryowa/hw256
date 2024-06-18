package service

import (
	"errors"
	"homework-1/internal/util"
	"testing"
)

func TestGetPackageType(t *testing.T) {
	film := NewFilm()
	packet := NewPacket()
	box := NewBox()

	if pt := GetPackageType(film); pt != string(filmType) {
		t.Errorf("Expected package type %v, got %v", filmType, pt)
	}
	if pt := GetPackageType(packet); pt != string(packetType) {
		t.Errorf("Expected package type %v, got %v", packetType, pt)
	}
	if pt := GetPackageType(box); pt != string(boxType) {
		t.Errorf("Expected package type %v, got %v", boxType, pt)
	}
}

func TestGetPackagePrice(t *testing.T) {
	film := NewFilm()
	packet := NewPacket()
	box := NewBox()

	if pp := GetPackagePrice(film); pp != float64(filmPrice) {
		t.Errorf("Expected package price %v, got %v", filmPrice, pp)
	}
	if pp := GetPackagePrice(packet); pp != float64(packetPrice) {
		t.Errorf("Expected package price %v, got %v", packetPrice, pp)
	}
	if pp := GetPackagePrice(box); pp != float64(boxPrice) {
		t.Errorf("Expected package price %v, got %v", boxPrice, pp)
	}
}

func TestFilmValidate(t *testing.T) {
	film := NewFilm()
	if err := film.Validate(100); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if err := film.Validate(0); !errors.Is(err, util.ErrWeightExceeds) {
		t.Errorf("Expected error %v, got %v", util.ErrWeightExceeds, err)
	}
}

func TestPacketValidate(t *testing.T) {
	packet := NewPacket()
	if err := packet.Validate(5); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if err := packet.Validate(10); !errors.Is(err, util.ErrWeightExceeds) {
		t.Errorf("Expected error %v, got %v", util.ErrWeightExceeds, err)
	}
}

func TestBoxValidate(t *testing.T) {
	box := NewBox()
	if err := box.Validate(20); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if err := box.Validate(30); !errors.Is(err, util.ErrWeightExceeds) {
		t.Errorf("Expected error %v, got %v", util.ErrWeightExceeds, err)
	}
}
