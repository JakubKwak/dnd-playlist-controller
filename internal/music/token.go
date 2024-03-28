package music

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

const tokenFile = "token.json"

func saveToken(token *oauth2.Token) error {
	jsonStr, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("couldn't convert token to json now wtf do we do: %w", err)
	}

	err = os.WriteFile(tokenFile, jsonStr, 0777)
	if err != nil {
		return fmt.Errorf("cringefest when writing file to my ass: %w", err)
	}

	return nil
}

func readToken() (*oauth2.Token, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(data, &token)

	return &token, err
}
