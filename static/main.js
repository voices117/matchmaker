function main(container) {
    let url = location.protocol.replace('http', 'ws') + location.host;
    let ws = new WebSocket(url + "/matchmaker", "matchmaker");

    ws.onmessage = (event) => {
        let data = JSON.parse(event.data);

        if (data.GameRoom !== undefined) {
            container.appendChild(createEventBox('Match ready!'));
            container.appendChild(createEventBox('Joining game...'));

            setTimeout(function() {
                // redirect to game room
                window.location.href = `/game.html?room_id=${data.GameRoom}&player_id=${getPlayerId()}`;
            }, 2000)

        } else {
            container.appendChild(createEventBox(data));
            container.scrollTop = container.scrollHeight;
        }
    }

    ws.onerror = function(error) {
        alert(`[error] ${error.message}`);
    }

    ws.onopen = () => {
        container.appendChild(createEventBox('Connected, waiting for match...'));

        // send login message
        ws.send(JSON.stringify({ client_id: getPlayerId() }))
    }

    ws.onclose = function(event) {
        var msg = '';
        if (event.wasClean) {
            msg = `Connection closed cleanly, code=${event.code} reason=${event.reason}`;
        } else {
            msg = 'Connection died';
        }
        container.appendChild(createEventBox(msg));
    }
}

function getPlayerId() {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    return params.player_id;
}


function createEventBox(msg) {
    var e = document.createElement('div');
    e.className = 'event';
    e.innerText = msg;
    return e
}