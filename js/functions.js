var chatinfo = '{"title":"Error"}';
var chatinfo_json = JSON.parse(chatinfo);
var user = JSON.parse('{"name":"Неизвестный"}');

function set_name() {
  localStorage.setItem(
    "name",
    prompt(
      'Добро пожаловать в "' +
        chatinfo_json.title +
        '"!\nВведите ваше имя, чтобы продолжить:'
    )
  );
}

function set_key() {
  localStorage.setItem(
    "key",
    prompt("Введите криптографический ключ, чтобы продолжить:")
  );
}

function encryptText(text, key) {
  const encrypted = CryptoJS.AES.encrypt(text, key).toString();
  return encrypted;
}

// Функция для расшифрования текста
function decryptText(encryptedText, key) {
  const decrypted = CryptoJS.AES.decrypt(encryptedText, key);
  const originalText = decrypted.toString(CryptoJS.enc.Utf8);
  return originalText;
}

function splitString(inputString, delimiter) {
  // Проверяем, что входная строка не пустая
  if (inputString.length === 0) {
    return [];
  }

  // Используем метод split для разделения строки
  return inputString.split(delimiter);
}
