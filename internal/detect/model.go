package detect

type DetectedStack struct {
	PHP     string `json:"php,omitempty"`
	Laravel string `json:"laravel,omitempty"`
	Nuxt    string `json:"nuxt,omitempty"`
	Vue     string `json:"vue,omitempty"`
	NuxtUI  string `json:"nuxt_ui,omitempty"`
}
