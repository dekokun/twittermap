const _ = require('underscore');

initialize = () => {
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

  // ゴールをおいておいてあげると親切かも
  const goalLatLng = new google.maps.LatLng(35.339126, 139.486817);
  const marker = new google.maps.Marker({
    position: goalLatLng,
    map: map,
    label: {
      fontWeight: 'bold',
      text: 'ゴール地点',
    }
  });

  let beforeTweets = [];
  let firstFlag = true;
  const onSuccess = (json) => {
    // _.differenceは同じオブジェクトかを比較するため使えないので自前で
    var difference = json.filter((obj) => { return !_.findWhere(beforeTweets, {url: obj.url}); });
    beforeTweets = json;
    const tweetsData = difference.map((tweet) => {
        const coordinates = tweet.coordinates;
        const lon = coordinates[0];
        const lat = coordinates[1];
        const time = new Date(tweet.created_at);
        const timeString = (time.getMonth()+1)+'/'+time.getDate()+' '+time.getHours()+':'+time.getMinutes();
        let text = '<dl><dt><a href="'+_.escape(tweet.url)+'" target="_blank">'+timeString+'</a></dt><dd>'+_.escape(decodeURIComponent(tweet.text))+'</dd></dl>';
        if (tweet.image_url) {
        text += '<a href="'+_.escape(tweet.image_url)+'" target="_blank"><img width="150" src="'+_.escape(tweet.image_url)+'" /></a>';
        }
        return {
          position: new google.maps.LatLng(lat, lon),
          text: text,
          image_url: tweet.image_url
          };
        });
    const markers = _.map(tweetsData, (tweet) => {
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
      map.panTo(_.last(tweetsData).position);
      firstFlag = false;
    }
    if (markers.length > 0) {
      const lastMarker = _.last(markers);
      // 新作ツイートが来たら内容を表示してあげましょう
      google.maps.event.trigger(lastMarker, 'click');
    }
  };

  var descriptions = [
    createContainer('ピン・アイコンを押すとツイートが表示されます。カメラアイコンは画像付きツイートです。ツイートは自動更新されます。'),
  ];
  _.each(descriptions, (container) => {
    map.controls[google.maps.ControlPosition.RIGHT_BOTTOM].push(container);
  });
  setTimeout(() => {
    map.controls[google.maps.ControlPosition.RIGHT_BOTTOM].clear();
  }, 10000);

  const main = async () => {
    const resp = await fetch("http://twittermap.dekokun.info/2018-03-17/tweets.json")
    if (resp.status >= 400) {
      alert('リクエスト失敗。作者にお問い合わせください');
    }
    const json = await resp.json();
    onSuccess(json);
  };
  setTimeout(main, 0.5 * 1000);
  setInterval(main, 60 * 1000);
};

google.maps.event.addDomListener(window, 'load', initialize);
