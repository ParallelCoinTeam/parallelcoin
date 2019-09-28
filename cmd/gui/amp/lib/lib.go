package lib

var AMP = `

`

///////////////

var AMPLIBS = map[string]string{
	"amp-boot":                  "amp-boot",
	"amp-timeago":               "amp-timeago",
	"amp-mathml":                "amp-mathml",
	"amp-form":                  "amp-form",
	"amp-iframe":                "amp-iframe",
	"amp-bind":                  "amp-bind",
	"amp-list":                  "amp-list",
	"amp-mustache":              "amp-mustache",
	"amp-sidebar":               "amp-sidebar",
	"amp-access":                "amp-access",
	"amp-analytics":             "amp-analytics",
	"amp-selector":              "amp-selector",
	"amp-animation":             "amp-animation",
	"amp-youtube":               "amp-youtube",
	"amp-carousel":              "amp-carousel",
	"amp-accordion":             "amp-accordion",
	"amp-twitter":               "amp-twitter",
	"amp-instagram":             "amp-instagram",
	"amp-lightbox":              "amp-lightbox",
	"amp-video":                 "amp-video",
	"amp-experiment":            "amp-experiment",
	"amp-live-list":             "amp-live-list",
	"amp-fit-text":              "amp-fit-text",
	"amp-viz-vega":              "amp-viz-vega",
	"amp-fx-collection":         "amp-fx-collection",
	"amp-social-share":          "amp-social-share",
	"amp-dynamic-css-classes":   "amp-dynamic-css-classes",
	"amp-position-observer":     "amp-position-observer",
	"amp-orientation-observer":  "amp-orientation-observer",
	"amp-install-serviceworker": "amp-install-serviceworker",
}

///////////////

func AMPlogo() string {
	return `<svg xmlns="http://www.w3.org/2000/svg" id="parallel" viewBox="0 0 108 128" width="108" height="128"><path class="logofill" d="M77.08,2.55c3.87,1.03 6.96,2.58 10.32,4.64c5.93,3.87 10.58,8.51 14.19,14.71c3.87,6.19 5.42,13.16 5.42,20.64c0,7.22 -1.81,14.18 -5.41,20.37c-3.61,6.45 -8.25,11.35 -14.19,14.96c-3.35,2.06 -6.96,3.87 -10.32,4.9c-3.87,1.03 -7.74,1.55 -11.61,1.55v-14.45c6.96,-0.26 13.42,-2.58 19.09,-8c5.67,-5.42 8.51,-11.87 8.51,-19.61c0,-7.74 -2.58,-14.19 -7.99,-19.6c-5.42,-5.42 -11.86,-8 -19.6,-8c-7.74,0 -14.44,2.58 -19.6,8c-5.42,5.42 -8,11.87 -8,19.6l0,85.9c-3.1,-3.1 -7.99,-7.74 -13.93,-13.67v-72.23c0,-3.87 0.52,-7.73 1.55,-11.35c1.03,-3.87 2.58,-7.22 4.64,-10.32c3.87,-5.93 8.52,-10.58 14.71,-14.45c6.19,-3.61 13.16,-5.16 20.64,-5.16c3.87,0 8,0.52 11.61,1.55zM78.37,42.28c0,7.22 -5.93,13.16 -13.15,13.16c-7.48,0.26 -13.16,-5.68 -13.16,-13.16c0,-7.22 5.94,-13.16 13.16,-13.16c7.22,0 13.15,5.93 13.15,13.16zM13.63,37.12l0,69.39c-6.19,-6.19 -11.09,-10.83 -13.93,-13.93l0,-55.46z" /></svg>`
}

func AMPlib() string {
	return `
<link rel="preload" as="script" href="https://cdn.ampproject.org/v0.js">
<script async src="https://cdn.ampproject.org/v0.js"></script>
<script async custom-element="amp-timeago" src="https://cdn.ampproject.org/v0/amp-timeago-0.1.js"></script>
<script async custom-element="amp-mathml" src="https://cdn.ampproject.org/v0/amp-mathml-0.1.js"></script>
<script async custom-element="amp-form" src="https://cdn.ampproject.org/v0/amp-form-0.1.js"></script>
<script async custom-element="amp-iframe" src="https://cdn.ampproject.org/v0/amp-iframe-0.1.js"></script>
<script async="" custom-element="amp-bind" src="https://cdn.ampproject.org/v0/amp-bind-0.1.js"></script>
<script async="" custom-element="amp-list" src="https://cdn.ampproject.org/v0/amp-list-0.1.js"></script>
<script async="" custom-template="amp-mustache" src="https://cdn.ampproject.org/v0/amp-mustache-0.2.js"></script>
<script async="" custom-element="amp-sidebar" src="https://cdn.ampproject.org/v0/amp-sidebar-0.1.js"></script>
<script async="" custom-element="amp-access" src="https://cdn.ampproject.org/v0/amp-access-0.1.js"></script>
<script async="" custom-element="amp-analytics" src="https://cdn.ampproject.org/v0/amp-analytics-0.1.js"></script>
<script async="" custom-element="amp-selector" src="https://cdn.ampproject.org/v0/amp-selector-0.1.js"></script>
<script async="" custom-element="amp-animation" src="https://cdn.ampproject.org/v0/amp-animation-0.1.js"></script>
<script async="" custom-element="amp-youtube" src="https://cdn.ampproject.org/v0/amp-youtube-0.1.js"></script>
<script async="" custom-element="amp-carousel" src="https://cdn.ampproject.org/v0/amp-carousel-0.1.js"></script>
<script async="" custom-element="amp-accordion" src="https://cdn.ampproject.org/v0/amp-accordion-0.1.js"></script>
<script async="" custom-element="amp-twitter" src="https://cdn.ampproject.org/v0/amp-twitter-0.1.js"></script>
<script async="" custom-element="amp-instagram" src="https://cdn.ampproject.org/v0/amp-instagram-0.1.js"></script>
<script async="" custom-element="amp-lightbox" src="https://cdn.ampproject.org/v0/amp-lightbox-0.1.js"></script>
<script async="" custom-element="amp-video" src="https://cdn.ampproject.org/v0/amp-video-0.1.js"></script>
<script async custom-element="amp-experiment" src="https://cdn.ampproject.org/v0/amp-experiment-0.1.js"></script>
<script async="" custom-element="amp-live-list" src="https://cdn.ampproject.org/v0/amp-live-list-0.1.js"></script>
<script async="" custom-element="amp-fit-text" src="https://cdn.ampproject.org/v0/amp-fit-text-0.1.js"></script>
<script async custom-element="amp-viz-vega" src="https://cdn.ampproject.org/v0/amp-viz-vega-0.1.js"></script>
<script async custom-element="amp-fx-collection" src="https://cdn.ampproject.org/v0/amp-fx-collection-0.1.js" ></script>
<script async="" custom-element="amp-social-share" src="https://cdn.ampproject.org/v0/amp-social-share-0.1.js"></script>
<script async custom-element="amp-dynamic-css-classes" src="https://cdn.ampproject.org/v0/amp-dynamic-css-classes-0.1.js"></script>
<script async="" custom-element="amp-position-observer" src="https://cdn.ampproject.org/v0/amp-position-observer-0.1.js"></script>
<script async="" custom-element="amp-orientation-observer" src="https://cdn.ampproject.org/v0/amp-orientation-observer-0.1.js"></script>
<script async custom-element="amp-install-serviceworker" src="https://cdn.ampproject.org/v0/amp-install-serviceworker-0.1.js"></script>


<style amp-boilerplate>body{-webkit-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-moz-animation:-amp-start 8s steps(1,end) 0s 1 normal both;-ms-animation:-amp-start 8s steps(1,end) 0s 1 normal both;animation:-amp-start 8s steps(1,end) 0s 1 normal both}@-webkit-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-moz-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-ms-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@-o-keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}@keyframes -amp-start{from{visibility:hidden}to{visibility:visible}}</style><noscript><style amp-boilerplate>body{-webkit-animation:none;-moz-animation:none;-ms-animation:none;animation:none}</style></noscript>
`
}

func AMPsw() string {
	return `<amp-install-serviceworker src="/sw/sw.js"
                               data-iframe-src="install-sw.html"
                               layout="nodisplay">
</amp-install-serviceworker>`
}
