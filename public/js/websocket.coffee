###
 * Websocket Client
 *
 * A simple websocket client
 *
 * @author Aotoki
###

###
 * Template Plugin
###

$.fn.tmpl = (d) ->
  s = $(this[0]).html().trim()
  if d
    for k, v of d
      s = s.replace(new RegExp('\\${' + k + '}', 'g'), v)
  return $(s)

###
 * Core
###

body = $("body, html")
messageArea = $("#message-area")
messageForm = $("#message-form")
message = $("#message")
inputFieldset = $("#input-fieldset")
onlineCount = $("#online-count")

messageTemplate = $("#message-template")

createMessage = (nickname, content) ->
  nickname = $("<p>").text(nickname).html()
  content = $("<p>").text(content).html()
  return messageTemplate.tmpl({nickname: nickname, message: content}).appendTo(messageArea)

updateOnlineCount = (count) ->
  onlineCount.text(count)

loadLogs = (packets) ->
  if packets
    for packet in packets
      handlePacket(packet)

handlePacket = (packet) ->
  switch packet.Type
    when 'message' then createMessage(packet.Data.Name, packet.Data.Content)
    when 'system' then createMessage("System", packet.Data).addClass("system-message")
    when 'online-count' then updateOnlineCount(packet.Data)
    when 'logs' then loadLogs(packet.Data)

  # Auto Scroll
  body.scrollTop(messageArea.height())

host = window.location.host
ws = new WebSocket("ws://" + host + "/chatroom")
ws.binaryType = "arraybuffer"
ws.onopen = (event)->
  inputFieldset.removeAttr("disabled")

ws.onclose = (event) ->
  createMessage("System", "You are lost connection, refresh page to reconnect.").addClass("system-message")
  inputFieldset.attr("disabled", "disabled")

ws.onmessage = (event)->
    packet = msgpack.unpack(new Uint8Array(event.data))
    handlePacket(packet)

messageForm.submit (e)->
  ws.send(message.val())
  message.val('')
  e.preventDefault()
