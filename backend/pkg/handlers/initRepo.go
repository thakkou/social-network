package handlers

import "01social/pkg/repository"

var Repos *repository.Repositories

func Init(r *repository.Repositories) {
	Repos = r
}
