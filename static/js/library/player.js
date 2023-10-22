// create youtube player
let player;
function onYouTubePlayerAPIReady() {
	player = new YT.Player('player', {
		height: '200',
		width: '300',
		playerVars: {
			'autoplay': 1,
			'start': 0
		},
		events: {
			'onReady': onPlayerReady,
			'onStateChange': onPlayerStateChange
		}
	});
}