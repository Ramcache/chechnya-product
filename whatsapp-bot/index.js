const express = require("express")
const {
    default: makeWASocket,
    useMultiFileAuthState,
    DisconnectReason,
    fetchLatestBaileysVersion
} = require("@whiskeysockets/baileys")
const qrcode = require("qrcode-terminal")
const fs = require("fs")

const app = express()
const port = 5555
app.use(express.json())

let sock
let latestQR = null // 📦 Храним последний QR

async function startBot() {
    const { version } = await fetchLatestBaileysVersion()
    const { state, saveCreds } = await useMultiFileAuthState("auth")

    sock = makeWASocket({
        version,
        auth: state,
        printQRInTerminal: false,
    })

    sock.ev.on("creds.update", saveCreds)

    sock.ev.on("connection.update", (update) => {
        const { connection, lastDisconnect, qr } = update

        if (qr) {
            latestQR = qr
            console.log("🔐 Сканируй QR-код ниже для подключения WhatsApp:")
            qrcode.generate(qr, { small: true })
        }

        if (connection === "close") {
            const shouldReconnect = lastDisconnect?.error?.output?.statusCode !== DisconnectReason.loggedOut
            console.log("❌ Соединение закрыто. Переподключение:", shouldReconnect)
            if (shouldReconnect) {
                startBot()
            }
        }

        if (connection === "open") {
            console.log("✅ Бот подключён к WhatsApp!")
        }
    })
}

// 📮 Endpoint для отправки кода
app.post("/send-code", async (req, res) => {
    const { phone, code } = req.body
    if (!phone || !code) {
        return res.status(400).send("Missing phone or code")
    }

    const chatId = phone.replace(/[^0-9]/g, "") + "@s.whatsapp.net"

    try {
        await sock.sendMessage(chatId, { text: `🔐 Ваш код подтверждения: ${code}` })
        console.log(`✅ Код ${code} отправлен на ${phone}`)
        res.send("Code sent")
    } catch (err) {
        console.error("❌ Ошибка при отправке:", err)
        res.status(500).send("Failed to send code")
    }
})

// 📤 Endpoint для получения QR-кода (текстом)
app.get("/qr", (req, res) => {
    if (!latestQR) {
        return res.status(404).send("QR-код пока не сгенерирован")
    }
    res.json({ qr: latestQR }) // можно отобразить через фронт или скопировать
})

startBot()
app.listen(port, () => {
    console.log(`🚀 WhatsApp Bot listening at http://localhost:${port}`)
})
