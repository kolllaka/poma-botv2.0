class Playlist {
	constructor(selector, songs = [], title, handler) {
		this.$el = document.getElementById(selector)

		this.song = ""
		this.songs = songs
		this.handler = handler
		this.title = title

		this.#render()
		this.#updateSongList()
		this.#setup()
	}

	#render() {
		this.$el.classList.add("video__playlist", "playlist")
		this.$el.innerHTML = getTemplate(this.title)
	}

	#updateSongList() {
		let musicTime = 0
		this.$el.querySelector(".playlist__body").innerHTML = this.songs.map((song, index) => {
			musicTime += song.duration
			return getTemplateItem(song, index)
		}).join("")
		this.$el.querySelector(".playlist__info").innerHTML = `общее количество треков: ${this.songs.length} на ${durationFormat(musicTime)}`
	}

	#setup() {
		this.clickHandler = this.clickHandler.bind(this)
		this.$el.addEventListener('click', this.clickHandler)
	}

	shuffle() {
		this.songs = this.songs.sort(() => Math.random() - .5)
		this.#updateSongList()
	}

	addSong(song) {
		this.songs.push(song)
		this.#updateSongList()
	}

	nextSong() {
		if (this.songs.length > 0) {
			const song = this.songs.shift()
			this.#updateSongList()
			this.song = song

			return song
		}

		this.song = ""
		console.log("playlist end");
	}

	clickHandler($event) {
		if ($event.target.classList.contains('del')) {
			let index = $event.target.dataset.index

			const song = this.songs.splice(index, 1)
			console.log(song[0]);
			this.#updateSongList()

			if (this.handler) {
				this.handler({ song: song[0], reason: "delete" })
			}

			return
		}
	}
}

const getTemplate = (title) => {
	return `
		<div class="playlist__title">${title}</div>
		<div class="playlist__info">${title}</div>
		<div class="playlist__item itemplaylist itemplaylist-title">
			<div class="itemplaylist__body">
				<div class="itemplaylist__cell">№№</div>
				<div class="itemplaylist__cell">кто заказал:</div>
				<div class="itemplaylist__cell itemplaylist__cell-n">название:</div>
				<div class="itemplaylist__cell">время</div>
				<div class="itemplaylist__cell">удалить</div>
			</div>
		</div>
		<ul class="playlist__body">
		</ul>
	`
}


const getTemplateItem = (song, index) => {
	const duration = durationFormat(song.duration)
	return `
	<!-- item -->
	<li class="playlist__item itemplaylist">
		<div class="itemplaylist__body">
			<div class="itemplaylist__cell">${index + 1}</div>
			<div class="itemplaylist__cell">${song.name || "owner"}</div>
			<div class="itemplaylist__cell itemplaylist__cell-n">${song.title}</div>
			<div class="itemplaylist__cell">${duration}</div>
			<div data-index="${index}" class="itemplaylist__cell del btn"></div>
		</div>
	</li>
	`
}


const durationFormat = (duration) => {
	let hours = (duration / 3600) | 0

	let minutes = ((duration % 3600) / 60) | 0

	let seconds = (duration % 60) | 0
	if (seconds < 10) seconds = `0` + seconds

	if (hours == 0) {
		return `${minutes}:${seconds}`
	}

	if (minutes < 10) minutes = `0` + minutes
	return `${hours}:${minutes}:${seconds}`
}