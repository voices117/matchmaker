function main(container) {
    let url = location.protocol.replace('http', 'ws') + location.host;
    let ws = new WebSocket(url + "/game", "game");
    let play = undefined;

    ws.onmessage = (event) => {
        container.innerHTML = "";

        let data = JSON.parse(event.data);
        if (data.play !== undefined) {
            play = data.play;
        } else if (data.error != "") {
            alert(data.error);
        } else {
            if (data.state != 'Unfinished') {
                let title = document.getElementById('title');
                var url = `/?player_id=${getPlayerId()}`;
                title.innerHTML = data.state + `<br><a href="${url}"> Back to lobby</a>`;
                container.appendChild(title);
            }

            container.appendChild(createBoard(data.board, ws));
        }
    }

    ws.onerror = function(error) {
        alert(`[error] ${error.message}`);
    }

    ws.onopen = () => {
        container.innerHTML += '<p>[Connected]</p>';

        // send login message
        // TODO: un-hardcode the client_id
        ws.send(JSON.stringify({ client_id: getPlayerId(), game_room_id: getGameRoom() }))
    }

    ws.onclose = function(event) {
        if (event.wasClean) {
            container.innerHTML += `<p>[close] Connection closed cleanly, code=${event.code} reason=${event.reason}</p>`;
        } else {
            container.innerHTML = '<p>[close] Connection died</p>';
        }
    }
}

function createBoard(board, ws) {
    let table = document.createElement('table');
    for (var i = 0; i < 3; i++) {
        let row = document.createElement('tr');
        for (var j = 0; j < 3; j++) {
            let cell = document.createElement('td');

            let value = document.createElement('a');
            value.href = '#';
            value.innerText = board[i * 3 + j];

            if (value.innerText == " ") {
                let index = i * 3 + j;
                value.onclick = function() {
                    ws.send(JSON.stringify({ position: index }));
                }
            }
            cell.appendChild(value);
            row.appendChild(cell);
        }

        table.appendChild(row);
    }
    return table;
}

// Read the 'room_id' from the URL query params
function getGameRoom() {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    return params.room_id;
}

function getPlayerId() {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    return params.player_id;
}