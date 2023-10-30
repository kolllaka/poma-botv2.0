# poma-bot v2.0
###
1.  В папку /static/aug/ помещаем все видео для гаданий
2.  в файле _.env_
    - USERID — ID twitch
    - YOUTUBEKEY — API key youtube
        > [где получить](https://console.cloud.google.com/apis/credentials) зарегистрироваться и во вкладке Сredentials получит API Key
    - AUDIOPATH — абсолютный путь до папки с музыкой
        > в папке будет подгружаться музыка с расширениями: .mp3, .mp4, .webm (остальные форматы игнорятся)
        > в названии файлов не должно быть _#_, треки с таким названием пропускаются
3. в файле _config.yaml_
    - rewardTitle — что отображается в награде
        > __ВАЖНО__ %s нужно, вместо него помещается прогой имя кто заказал награду
        > 
        > (пример): '%s сегодня'
    - rewardName — имя награды (регистр не важен)
        > (пример): 'Гадание'
    - duration — максимальное время заказа в __секундах__ (если время заказа выше этого он не добавится в очередь)
        > только к аудио заказам, по умолчанию 600
        > 
        > (пример): 600
4. запускаем _poma-botv2.0.exe_ и на:
    - http://localhost:8080/aug — гадания (подключаем в ОБС)
    - http://localhost:8080/music — очередь заказов (открываем в браузере)
