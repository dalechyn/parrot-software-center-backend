package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"parrot-software-center-backend/tokens"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	log.Debug("RefreshToken request attempt")

	// Decoding incoming token renewal request
	inRequest := &tokenRequest{}
	if err := json.NewDecoder(r.Body).Decode(inRequest); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	// If found, generating new token pair
	newAccessToken, newRefreshToken, err := tokens.UpdateTokens(inRequest.RefreshToken)
	if err != nil {
		log.Error(err)
		http.Error(w, "RefreshToken pair update failed", http.StatusInternalServerError)
	}

	// JSON encoding response with new token pair
	resBytes, err := json.Marshal(&tokenResponse{newAccessToken, newRefreshToken})
	if err != nil {
		log.Error(err)
	}

	// http Response
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resBytes); err != nil {
		log.Error(err)
	}
}
