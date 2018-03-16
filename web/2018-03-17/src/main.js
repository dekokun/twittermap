const escape = require('escape-html');

google.maps.event.addDomListener(window, 'load', () => {
  const zIndexDefault = 0;
  const zIndexHigh = 100;
  const mapOptions = {
    // 大体道中の中心くらい
    center: { lat: 35.529344, lng: 139.643372},
    zoom: 10
  };
  const map = new google.maps.Map(document.getElementById('map-canvas'),
      mapOptions);

  const createContainer = (message) => {
    const container = document.createElement('div');
    container.style.margin = '10px';
    container.style.padding = '10px';
    container.style.border = '0px solid #000';
    container.style.background = '#FFF';
    container.innerHTML = message;
    return container;
  }
  let globalInfoWindow;
  const addMarker = (position, text, type, zIndex) => {
    let markerUrl;
    switch (type) {
      case 'comment':
        markerUrl = 'https://maps.google.com/mapfiles/ms/icons/pink-dot.png';
        break;
      case 'start':
        markerUrl = 'img/start.png';
        break;
      case 'camera':
        markerUrl = 'img/camera.png';
        break;
      default:
        markerUrl = 'https://maps.google.com/mapfiles/ms/icons/red-dot.png';
    }
    const marker = new google.maps.Marker({
      position: position,
      map: map,
      animation: google.maps.Animation.DROP,
      zIndex : zIndex,
      icon: markerUrl
    });
    const myInfoWindow = new google.maps.InfoWindow({
      content: text,
      disableAutoPan: true,
      zIndex : zIndex,
      maxWidth: 150
    });
    google.maps.event.addListener(marker, 'click', () => {
      if (globalInfoWindow && globalInfoWindow.close) {
        globalInfoWindow.close();
      }
      globalInfoWindow = myInfoWindow;
      globalInfoWindow.open(map,marker);
    });
    return marker;
  };

  // ゴールとスタートをわかりやすくしておく
  const startLatLng = new google.maps.LatLng(35.685527, 139.780991);
  new google.maps.Marker({
    position: startLatLng,
    map: map,
    label: {
      fontWeight: 'bold',
      text: 'スタート',
    }
  });
  const goalLatLng = new google.maps.LatLng(35.339126, 139.486817);
  new google.maps.Marker({
    position: goalLatLng,
    map: map,
    label: {
      fontWeight: 'bold',
      text: 'ゴール',
    }
  });

  let beforeTweets = [];
  let firstFlag = true;
  const onSuccess = (json) => {
    const difference = json.filter((obj1) => { return !beforeTweets.some((obj2) => {return obj1.id == obj2.id})});
    beforeTweets = json;
    const tweetsData = difference.map((tweet) => {
        const coordinates = tweet.coordinates;
        const lon = coordinates[0];
        const lat = coordinates[1];
        const time = new Date(tweet.created_at);
        const timeString = ("00" + (time.getMonth()+1)).slice(-2)+'/'+("00" + time.getDate()).slice(-2)+' '+("00" + time.getHours()).slice(-2)+':'+("00" + time.getMinutes()).slice(-2);
        let text = '<dl><dt><a href="'+escape(tweet.url)+'" target="_blank">'+timeString+'</a></dt><dd>'+escape(decodeURIComponent(tweet.text))+'</dd></dl>';
        if (tweet.image_url) {
        text += '<a href="'+escape(tweet.image_url)+'" target="_blank"><img width="150" src="'+escape(tweet.image_url)+'" /></a>';
        }
        return {
          position: new google.maps.LatLng(lat, lon),
          text: text,
          image_url: tweet.image_url
          };
        });
    const markers = tweetsData.map((tweet) => {
        let type, zIndex;
        if (tweet.text.match(/今日のスタート地点は/)) {
          type = 'start';
          zIndex = zIndexHigh;
        } else if (tweet.image_url) {
          type = 'camera';
          zIndex = zIndexDefault;
        } else {
          type = 'comment';
          zIndex = zIndexDefault;
        }
        return addMarker(tweet.position, tweet.text, type, zIndex);
        });
    // 最初の一回だけ、最後のピンの場所に中心をずらす
    if (tweetsData.length > 0 && firstFlag) {
      const [lastTweet] = tweetsData.slice(-1)
      map.panTo(lastTweet.position);
      firstFlag = false;
    }
    if (markers.length > 0) {
      const [lastMarker] = markers.slice(-1);
      // 新作ツイートが来たら内容を表示してあげましょう
      google.maps.event.trigger(lastMarker, 'click');
    }
  };

  var descriptions = [
    createContainer('ピン・アイコンを押すとツイートが表示されます。カメラアイコンは画像付きツイートです。ツイートは自動更新されます。'),
  ];
  descriptions.forEach( (container) => {
    map.controls[google.maps.ControlPosition.RIGHT_BOTTOM].push(container);
  });
  setTimeout(() => {
    map.controls[google.maps.ControlPosition.RIGHT_BOTTOM].clear();
  }, 10000);

  const main = async () => {
    const resp = await fetch("//twittermap.dekokun.info/2018-03-17/tweets.json")
    if (resp.status >= 400) {
      alert('リクエスト失敗。作者にお問い合わせください');
    }
    const json = await resp.json();
    onSuccess(json);
  };
  setTimeout(main, 0.5 * 1000);
  setInterval(main, 60 * 1000);
});
