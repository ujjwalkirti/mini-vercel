package client

import "github.com/supabase-community/supabase-go"

func NewSupabaseClient(API_URL string, API_KEY string) (*supabase.Client, error) {
	client, err := supabase.NewClient(API_URL, API_KEY, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
