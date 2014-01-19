/**
 * Websocket Client
 *
 * A simple websocket client
 *
 * @author Aotoki
 */

/**
 * Template Plugin
 */

$.fn.tmpl = function(d) {
  var s = $(this[0]).html().trim();
  if (d) {
    for (k in d) {
      s = s.replace(new RegExp('\\${' + k + '}', 'g'), d[k]);
    }
  }
  return $(s);
};

/**
 * Core
 */

var body = $("body, html");
var messageArea = $("#message-area");
var messageForm = $("#message-form");
var message = $("#message");
var inputFieldset = $("#input-fieldset");

var messageTemplate = $("#message-template");

function createMessage(nickname, content) {
  nickname = $("<p>").text(nickname).html();
  content = $("<p>").text(content).html();
  return messageTemplate.tmpl({nickname: nickname, message: content}).appendTo(messageArea);
}

var scrollHeight = 0;
function appendMessage(rawMessage) {
  var message = rawMessage.split(":");
  if(message.length == 2) { // Chat Message
    createMessage(message[0], message[1]);
  } else { // System Message
    createMessage("System", rawMessage).addClass("system-message");
  }

  // Auto Scroll
  body.scrollTop(messageArea.height());
}

var ws = new WebSocket("ws://localhost:3000/chatroom");
ws.onopen = function(event) {
  inputFieldset.removeAttr("disabled")
}
ws.onmessage = function(event) {
  appendMessage(event.data);
}

messageForm.submit(function(e) {
  ws.send(message.val());
  message.val('');
  e.preventDefault()
});
