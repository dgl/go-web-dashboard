let ws;

function show(name) {
  document.cookie = "gww_name=" + name +
    "; expires=" + new Date(new Date().getTime() + (5 * 365 * 24 * 60 * 60 * 1000)).toUTCString();

  if (typeof(Storage) !== undefined) {
    let url = window.localStorage.getItem('lasturl');
    if (url) {
      displayURL(url);
    } else {
      displayURL("/time");
    }
  }

  ws = connect(location.host, name);
  setInterval(function() {
    if (ws == null || ws.readyState != ws.OPEN) {
      ws = connect(location.host, name);
    } else {
      ws.send(JSON.stringify({type: "ping"}));
    }
  }, 60 * 1000);
}

function connect(host, name) {
  let ws = new WebSocket("ws://" + host + "/ws");

  ws.onopen = function() {
    ws.send(JSON.stringify({type: "connect", name: name}));
  }

  ws.onmessage = function(e) {
    let data = JSON.parse(e.data);
    console.log(data);

    if ('url' in data) {
      displayURL(data.url);
      if (typeof(Storage) !== undefined) {
        window.localStorage.setItem('lasturl', data.url);
      }
    }
  }

  return ws;
}

function displayURL(url) {
  let iframe = document.querySelector("#content");
  iframe.src = url;
}
