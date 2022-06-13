function main(container) {
    let url = location.protocol.replace('http', 'ws') + location.host;
    let ws = new WebSocket(url + "/matchmaker", "matchmaker");

    ws.onmessage = (event) => {
        let msg = document.createElement('p');
        msg.innerText = event.data;

        container.appendChild(msg);
    }

    ws.onerror = function(error) {
        alert(`[error] ${error.message}`);
    }

    ws.onopen = () => {
        container.innerHTML += '<p>[Connected]</p>';

        // send login message
        ws.send(JSON.stringify({ client_id: 'test' }))
    }

    ws.onclose = function(event) {
        if (event.wasClean) {
            container.innerHTML += `<p>[close] Connection closed cleanly, code=${event.code} reason=${event.reason}</p>`;
        } else {
            container.innerHTML = '<p>[close] Connection died</p>';
        }
    }
}