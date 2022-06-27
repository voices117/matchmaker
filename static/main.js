function main(container) {
    let url = location.protocol.replace('http', 'ws') + location.host;
    let ws = new WebSocket(url + "/matchmaker", "matchmaker");

    ws.onmessage = (event) => {
        let msg = document.createElement('p');
        msg.innerText = event.data;

        container.appendChild(msg);

        let data = JSON.parse(event.data);

        if (data.GameRoom !== undefined) {
            window.location.href = "/game.html?room_id=" + data.GameRoom;
        }
    }

    ws.onerror = function(error) {
        alert(`[error] ${error.message}`);
    }

    ws.onopen = () => {
        container.innerHTML += '<p>[Connected]</p>';

        // send login message
        ws.send(JSON.stringify({ client_id: getPlayerId() }))
    }

    ws.onclose = function(event) {
        if (event.wasClean) {
            container.innerHTML += `<p>[close] Connection closed cleanly, code=${event.code} reason=${event.reason}</p>`;
        } else {
            container.innerHTML = '<p>[close] Connection died</p>';
        }
    }
}

function getPlayerId() {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    return params.player_id;
}