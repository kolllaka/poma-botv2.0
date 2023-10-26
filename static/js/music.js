const constVolumeLevel = 15
let VolumeLevel = localStorage.getItem('VolumeLevel') || constVolumeLevel;
localStorage.setItem('VolumeLevel', VolumeLevel)

const sendSocket = (data) => {
	socket.send(JSON.stringify(data))
}

// video
const myPlayer = document.getElementById("myplayer")

const music = new Playlist("music", [], "очередь заказов", sendSocket)

const myMusic = new Playlist("mymusic", [], "мой плейлист")


let playlistStruct = []
const audio = new Audio()
const getMetaData = () => {
	const song = playlistStruct.shift()


	if (song) {
		console.log("shift", song.link);
		audio.src = song.link
		audio.onloadedmetadata = () => {
			song.duration = audio.duration

			//localStorage.setItem(song.name, song.duration)
			myMusic.addSong(song);
			audio.pause()
			audio.src = ""

			getMetaData()
		}

		audio.onerror = () => {
			getMetaData()
		}
	}
}

let isResizeble = false
const WEBSOCKET = 'ws://127.0.0.1:8080/music/ws'
const handler = () => {
	//console.log(msgStruct)
	if (msgStruct.isreward) {
		music.addSong(msgStruct)

		return
	}

	// const duration = localStorage.getItem(msgStruct.name)
	// if (duration) {
	// 	msgStruct.duration = +duration

	// 	console.log("from cashe");
	// 	myMusic.addSong(msgStruct);

	// 	return
	// }
	console.log("msgStruct.name", msgStruct.link);

	playlistStruct.push(msgStruct)

	if (!isResizeble) {
		getMetaData()

		isResizeble = true
	}
}

connectWS(WEBSOCKET, handler)

// volume 
const volume = document.getElementById('volume')
volume.value = VolumeLevel
const setMasterVolume = (volumeVal) => {
	player.setVolume(volumeVal);
	myPlayer.volume = volumeVal / 100
}
volume.addEventListener("input", (e) => {
	localStorage.setItem('VolumeLevel', e.target.value)
	setMasterVolume(e.target.value)
});

// autoplay video
function onPlayerReady(event) {
	player.loadVideoById("QxtKHo0iMa4");
	setMasterVolume(VolumeLevel)

	event.target.playVideo()
}



let isYPlaying = false,
	isPlaying = false,
	isYPlay = isYPlaying
const toggleBtn = document.getElementById("playerControl").querySelector("[data-btn='play']")
console.log(toggleBtn);
const togglePlay = () => {
	if (isYPlaying || isPlaying) {
		toggleBtn.classList.add("btn-play")
		toggleBtn.classList.remove("btn-pause")
	} else {
		toggleBtn.classList.remove("btn-play")
		toggleBtn.classList.add("btn-pause")
	}
}
// when video change state
function onPlayerStateChange(event) {
	isYPlaying = false
	if (event.data == YT.PlayerState.ENDED) {
		nextSongHandler()

		return
	}

	if (event.data == YT.PlayerState.PLAYING) {
		isYPlaying = true
	}

	togglePlay()
}

myPlayer.onplay = (e) => {
	isPlaying = true
	togglePlay()
}

myPlayer.onpause = (e) => {
	isPlaying = false
	togglePlay()
}


// myPlayer
myPlayer.addEventListener('ended', (e) => {
	nextSongHandler()
})

const playedVideoTitle = document.getElementById("video__played")
const playedVideoOwned = document.getElementById("video__owned")
const nextSongHandler = () => {
	let song = music.nextSong()
	myPlayer.pause()
	myPlayer.src = ""

	if (song) {
		console.log(song)
		playedVideoTitle.innerText = song.title
		playedVideoOwned.innerText = song.name
		player.loadVideoById(song.link);


		sendSocket({ song: song, reason: "played" });

		return
	}

	player.stopVideo()
	song = myMusic.nextSong()
	if (song) {
		console.log(song)
		playedVideoTitle.innerText = song.title
		playedVideoOwned.innerText = "из моего плейлиста"
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
		case "play":
			if (isYPlaying || isPlaying) {
				// play
				isYPlay = isYPlaying
				myPlayer.pause()
				player.pauseVideo()
			} else {
				// pause
				if (isYPlay) {
					player.playVideo()
				} else {
					myPlayer.play()
				}
			}


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