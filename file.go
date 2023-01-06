package main

import (
	"fmt"
	"net/http"
	"strings"
)

func downloadFile() {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		// Générez dynamiquement le contenu du fichier
		var content strings.Builder
		for i := 1; i <= 10; i++ {
			content.WriteString(fmt.Sprintf("Ligne %d\n", i))
		}

		// Définissez les en-têtes Content-Type et Content-Disposition
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", `attachment; filename="pairsKraken.csv"`)

		// Écrivez le contenu du fichier dans le corps de la réponse
		w.Write([]byte(content.String()))

	})

	// Démarrez le serveur HTTP et écoutez sur le port 8080.
	http.ListenAndServe(":8080", nil)
}