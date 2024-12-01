document.addEventListener("DOMContentLoaded", async () => {
  var name_inp = document.getElementById("name");
  var key_inp = document.getElementById("key");
  var msg_inp = document.getElementById("msg");

  const new_message_audio = document.getElementById("new_message");
  new_message_audio.volume = 0.1;

  name_inp.addEventListener("input", () => {
    localStorage.setItem("name", name_inp.value);
  });

  key_inp.addEventListener("input", () => {
    localStorage.setItem("key", key_inp.value);
  });

  await fetch("/settings.json")
    .then((response) => {
      if (!response.ok) {
        throw new Error("Сеть ответила с ошибкой: " + response.status);
      }
      return response.text();
    })
    .then((data) => {
      chatinfo = data;
    })
    .catch((error) => {
      console.error("Ошибка при получении текста:", error);
    });

  chatinfo_json = JSON.parse(chatinfo);

  document.title = chatinfo_json.title + " - Ciphira message transfer protocol";
  document.getElementById("logo").innerHTML += " : " + window.location.host;

  var ws = new WebSocket("ws://" + window.location.host + "/ws");

  ws.onopen = function () {
    document.getElementById("logo").innerHTML +=
      " - Соединение с WebSocket-сервером установлено";
  };

  ws.onmessage = function (event) {
    new_message_audio.play();

    var messages = splitString(event.data, "\n");
    for (let index = 0; index < messages.length; index++) {
      const element = messages[index];
      document.getElementById("chat").innerHTML +=
        "<p></p>" + decryptText(element, localStorage.getItem("key"));
    }

    const scrollableDiv = document.getElementById("chat");
    scrollableDiv.scrollTop = scrollableDiv.scrollHeight;
  };

  ws.onclose = function () {
    console.log("WebSocket connection closed");
    document.getElementById("logo").innerHTML =
      " Ошибка: соединение с WebSocket-сервером разорвано";

    document.documentElement.style.setProperty("--primary", "red");
  };

  msg_inp.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      ws.send(
        encryptText(
          "(" +
            localStorage.getItem("name") +
            ") [" +
            new Date() +
            "] " +
            msg_inp.value +
            "",
          localStorage.getItem("key")
        )
      );

      msg_inp.value = "";
    }
  });

  if (!localStorage.getItem("hasVisited")) {
    set_name();
    set_key();

    name_inp.value = localStorage.getItem("name");

    localStorage.setItem("hasVisited", "true");
  } else {
  }

  name_inp.value = localStorage.getItem("name");
  key_inp.value = localStorage.getItem("key");

  // Пример использования
  const key = "my-secret-key"; // Ключ для шифрования
  const textToEncrypt = "Привет, мир!";
});
