package main

import (
	"fmt"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
	"io"
	"encoding/base64"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}


	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)
	file, _, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting thumbnail", err)
		return
	}

	defer file.Close()

	//contentType := header.Header.Get("Content-Type")

	img, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading file", err)
		return
	}

	md, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get video", err)
		return
	}
	if md.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Video does not belong to user", err)
		return
	}

	base_img := base64.StdEncoding.EncodeToString(img)
	url := fmt.Sprintf("data:text/plain;base64,%s", base_img)
	md.ThumbnailURL = &url

	err = cfg.db.UpdateVideo(md)
	if err != nil {
		//delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}


	respondWithJSON(w, http.StatusOK, md)
}
