// Bashayr Alabdullah
package main
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "YOUR DEVELOPER KEY"
const slackAppSecret = "YOUR_SIGNING_SECRET_HERE"

func handleRequestAndResponse(w http.ResponseWriter, r *http.Request) {
	//slack
	verifier, err := slack.NewSecretsVerifier(r.Header, slackAppSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// youtube
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
	call := service.Search.List("id,snippet").
		Q(s.Text).
		MaxResults(20)
	response, err := call.Do()
	log.Print(err, "")
	videos := make(map[string]string)
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		}
	}
	buildURLLink(w, "Videos", videos)
}

func main() {

	fmt.Println("[INFO] Server listening")
	http.HandleFunc("/videos", handleRequestAndResponse)

	_ = http.ListenAndServe(":8080", nil)
}

func buildURLLink(w http.ResponseWriter, sectionName string, matches map[string]string) {
	_, _ = fmt.Fprint(w, "\n", sectionName)
	for id, title := range matches {
		_, _ = fmt.Fprint(w, "\nhttps://www.youtube.com/watch?v=", id, " ", title, "\n")
	}
}
