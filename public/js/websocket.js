// Generated by CoffeeScript 1.6.3
/*
 * Websocket Client
 *
 * A simple websocket client
 *
 * @author Aotoki
*/


/*
 * Template Plugin
*/


(function() {
  var body, createMessage, handlePacket, host, inputFieldset, loadLogs, message, messageArea, messageForm, messageTemplate, onlineCount, updateOnlineCount, ws;

  $.fn.tmpl = function(d) {
    var k, s, v;
    s = $(this[0]).html().trim();
    if (d) {
      for (k in d) {
        v = d[k];
        s = s.replace(new RegExp('\\${' + k + '}', 'g'), v);
      }
    }
    return $(s);
  };

  /*
   * Core
  */


  body = $("body, html");

  messageArea = $("#message-area");

  messageForm = $("#message-form");

  message = $("#message");

  inputFieldset = $("#input-fieldset");

  onlineCount = $("#online-count");

  messageTemplate = $("#message-template");

  createMessage = function(nickname, content) {
    nickname = $("<p>").text(nickname).html();
    content = $("<p>").text(content).html();
    return messageTemplate.tmpl({
      nickname: nickname,
      message: content
    }).appendTo(messageArea);
  };

  updateOnlineCount = function(count) {
    return onlineCount.text(count);
  };

  loadLogs = function(packets) {
    var packet, _i, _len, _results;
    if (packets) {
      _results = [];
      for (_i = 0, _len = packets.length; _i < _len; _i++) {
        packet = packets[_i];
        _results.push(handlePacket(packet));
      }
      return _results;
    }
  };

  handlePacket = function(packet) {
    switch (packet.Type) {
      case 'message':
        createMessage(packet.Data.Name, packet.Data.Content);
        break;
      case 'system':
        createMessage("System", packet.Data).addClass("system-message");
        break;
      case 'online-count':
        updateOnlineCount(packet.Data);
        break;
      case 'logs':
        loadLogs(packet.Data);
    }
    return body.scrollTop(messageArea.height());
  };

  host = window.location.host;

  ws = new WebSocket("ws://" + host + "/chatroom");

  ws.binaryType = "arraybuffer";

  ws.onopen = function(event) {
    return inputFieldset.removeAttr("disabled");
  };

  ws.onclose = function(event) {
    createMessage("System", "You are lost connection, refresh page to reconnect.").addClass("system-message");
    return inputFieldset.attr("disabled", "disabled");
  };

  ws.onmessage = function(event) {
    var packet;
    packet = msgpack.unpack(new Uint8Array(event.data));
    return handlePacket(packet);
  };

  messageForm.submit(function(e) {
    ws.send(message.val());
    message.val('');
    return e.preventDefault();
  });

}).call(this);
