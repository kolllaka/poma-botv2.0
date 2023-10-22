const sendSocket = (data) => {
	socket.send(JSON.stringify(data))
}

// video
const myPlayer = document.getElementById("myplayer")


const music = new Playlist("music", [
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: 'hxqW1Eq0iP4', duration: 238 },
	// { name: 'KoLLIaka', title: '◤Nightcore◢ ↬ Criminal [Switching Vocals]', link: 'WCL-sYe02Xg', duration: 242 },
	// { name: 'KoLLIaka', title: 'Take Me On!', link: 'xmJGiUTSh5E', duration: 255 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: 'hxqW1Eq0iP4', duration: 238 },
	// { name: 'KoLLIaka', title: '◤Nightcore◢ ↬ Criminal [Switching Vocals]', link: 'WCL-sYe02Xg', duration: 242 },
	// { name: 'KoLLIaka', title: 'Take Me On!', link: 'xmJGiUTSh5E', duration: 255 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: 'hxqW1Eq0iP4', duration: 238 },
	// { name: 'KoLLIaka', title: '◤Nightcore◢ ↬ Criminal [Switching Vocals]', link: 'WCL-sYe02Xg', duration: 242 },
	// { name: 'KoLLIaka', title: 'Take Me On!', link: 'xmJGiUTSh5E', duration: 255 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: 'hxqW1Eq0iP4', duration: 238 },
	// { name: 'KoLLIaka', title: '◤Nightcore◢ ↬ Criminal [Switching Vocals]', link: 'WCL-sYe02Xg', duration: 242 },
	// { name: 'KoLLIaka', title: 'Take Me On!', link: 'xmJGiUTSh5E', duration: 255 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: 'hxqW1Eq0iP4', duration: 238 },
	// { name: 'KoLLIaka', title: '◤Nightcore◢ ↬ Criminal [Switching Vocals]', link: 'WCL-sYe02Xg', duration: 242 },
	// { name: 'KoLLIaka', title: 'Take Me On!', link: 'xmJGiUTSh5E', duration: 255 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: 'hxqW1Eq0iP4', duration: 238 },
	// { name: 'KoLLIaka', title: '◤Nightcore◢ ↬ Criminal [Switching Vocals]', link: 'WCL-sYe02Xg', duration: 242 },
	{ name: 'KoLLIaka', title: 'Take Me On!', link: 'xmJGiUTSh5E', duration: 255 },
], "очередь заказов", sendSocket)


const myMusic = new Playlist("mymusic", [
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: './static/playlist/3.mp3', duration: 238 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: './static/playlist/2.webm', duration: 238 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: './static/playlist/1.webm', duration: 238 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: './static/playlist/2.webm', duration: 238 },
	// { name: 'KoLLIaka', title: 'Courtesy Call - Thousand Foot Krutch (Lyrics)', link: './static/playlist/1.webm', duration: 238 },
], "мой плейлист")



const getMetaData = async (song) => {
	const audio = new Audio(song.link)
	audio.onloadedmetadata = () => {
		song.duration = audio.duration
		console.log(audio.audioTracks);
		myMusic.addSong(song);
	}
}

toServer("http://localhost:8080/api/playlist", "GET")
	.then((resp) => {
		console.log(resp);
		resp.forEach((song) => {
			getMetaData(song)

		})
	})
	.catch((err) => {
		console.log(err);
	})


const WEBSOCKET = 'ws://127.0.0.1:8080/music/ws'
const handler = () => {
	console.log(msgStruct)
	if (msgStruct.isreward) {
		music.addSong(msgStruct)

		return
	}

	myMusic.addSong(msgStruct)
}



connectWS(WEBSOCKET, handler)
// volume 
const volume = document.getElementById('volume')
const setMasterVolume = (volumeVal) => {
	player.setVolume(volumeVal);
	myPlayer.volume = volumeVal / 100
}
volume.addEventListener("input", (e) => setMasterVolume(e.target.value));

// autoplay video
function onPlayerReady(event) {
	player.loadVideoById("QxtKHo0iMa4");
	setMasterVolume(50)

	event.target.playVideo()
}

// when video ends
function onPlayerStateChange(event) {
	if (event.data == YT.PlayerState.ENDED) {
		nextSongHandler()
	}
}

// myPlayer
myPlayer.addEventListener('ended', (e) => {
	nextSongHandler()
})

const nextSongHandler = () => {
	let song = music.nextSong()
	myPlayer.pause()
	myPlayer.src = ""
	if (song) {
		console.log(song)
		player.loadVideoById(song.link);

		sendSocket({ song: song, reason: "played" });

		return
	}

	player.stopVideo()
	song = myMusic.nextSong()
	if (song) {
		console.log(song)
		myPlayer.src = song.link
		myPlayer.play()
	}
}

document.getElementById("playerControl").addEventListener("click", (e) => {
	e.preventDefault()
	switch (e.target.dataset.btn) {
		case "skip":
			skipBtnHandler(e)

			break
		case "send":
			sendBtnHandler(e)

			break
		case "shuffle":
			shuffledBtnHandler(e)

			break
	}

})

const skipBtnHandler = (e) => {
	if (music.song != "") {
		sendSocket({ song: music.song, reason: "skip" });
	}
	nextSongHandler()
}

const sendBtnHandler = (e) => {
	const url = "http://localhost:8080/api/yplaylist"
	const label = e.target.closest("label")
	const value = label.querySelector('input').value
	const reg = new RegExp(`list=([^&\"\']*)`)

	let result = value.match(reg)
	if (result) {
		console.log(result[1]);
		toServer(url, "POST", { link: result[1] })
			.then((response) => {
				console.log(response);
				response.forEach((song) => {
					myMusic.addSong(song)
				})
			})
			.catch((err) => {
				console.log("[error]", err);
			})
	}
}

const shuffledBtnHandler = (e) => {
	myMusic.shuffle()
}

