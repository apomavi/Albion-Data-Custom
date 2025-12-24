🛡️ Albion Online Guild Bank Monitor (Custom Client)
Bu proje, Albion Online lonca (guild) bankası hareketlerini ağ trafiği üzerinden dinleyerek, belirli bir limitin üzerindeki para çıkışlarını (Withdraw) tespit eder ve Discord üzerinden anlık bildirim gönderir.

Lonca ekonomisini korumak ve şüpheli işlemleri (RMT, hırsızlık vb.) anında fark etmek amacıyla geliştirilmiştir.

A quick note on the legality of this application and if it
violates the Terms and Conditions for Albion Online. Here is
the response from SBI when asked if we are allowed to do
monitor network packets relating to Albion Online:
> Our position is quite simple. As long as you just look and
analyze we are ok with it. The moment you modify or manipulate
something or somehow interfere with our services we will react
(e.g. perma-ban, take legal action, whatever).

~ MadDave - Technical Lead for Albion Online

Source: https://forum.albiononline.com/index.php/Thread/51604-Is-it-allowed-to-scan-your-internet-trafic-and-pick-up-logs/?postID=512670#post512670

This client monitors local network traffic, identifies UDP packets
that contain relevant data for Albion Online, and ships the information
off to a central NATS server that anyone can subscribe to.

🔗 Atıf ve Teşekkür (Credits)
Bu proje, harika bir açık kaynak projesi olan Albion Data Client altyapısı üzerine inşa edilmiştir.

Orijinal proje, Albion Online ağ paketlerini dinlemek ve ayrıştırmak (sniffing & parsing) için gerekli olan temel kütüphaneyi sağlar. Biz bu güçlü altyapıyı kullanarak, özellikle Guild Log paketlerine odaklanan ve bunları Discord Webhook ile entegre eden özelleştirilmiş bir versiyon geliştirdik.

Orijinal projeye buradan ulaşabilirsiniz: github.com/ao-data/albiondata-client

🚀 Özellikler
Paket Analizi: Oyunun ağ trafiğini dinler ve Guild işlem loglarını yakalar.

Akıllı Filtreleme: Sadece "Para Çekme" (Withdraw) işlemlerini filtreler.

Eşik Değeri (Threshold): Belirlenen miktar (örneğin 10 Milyon Silver) üzerindeki işlemler için alarm üretir.

Discord Entegrasyonu: Şüpheli işlemleri detaylı (Oyuncu adı, Miktar, Tarih) bir şekilde Discord kanalına raporlar.

Deduplication: Aynı logun tekrar tekrar gönderilmesini önlemek için son işlenen log zamanını (cursor) kaydeder.

🛠️ Kurulum ve Gereksinimler
Ön Hazırlık
Bu yazılımın çalışabilmesi için bilgisayarınızda ağ paketlerini yakalayacak bir sürücüye ihtiyaç vardır:

Windows: Npcap (Kurarken "Install Npcap in WinPcap API-compatible Mode" seçeneğini işaretleyin).

Linux/macOS: libpcap kütüphanesi.

Derleme (Build)
Projeyi bilgisayarınıza klonlayın ve proje dizininde şu komutu çalıştırın:



⚙️ Yapılandırma
Kod içerisindeki client paketinde bulunan ayarları kendi sunucunuza göre düzenlemelisiniz:

**To post to a specific Discord channel, create a webhook in that channel's editor and paste the link here: client/operation_guild_log.go ----> const DiscordWebhookURL = "Link_in_here"**

// Discord Webhook URL'nizi buraya girin
const DiscordWebhookURL = "https://discord.com/api/webhooks/..."

// Bildirim için alt limit (Örn: 10 Milyon Silver)
// Negatif değer girilmelidir (Para çıkışı olduğu için)
const NotificationThreshold = -10000000 
📸 Ekran Görüntüleri / Örnek Çıktı
Discord Bildirimi:

<img width="1292" height="882" alt="sss" src="https://github.com/user-attachments/assets/c5ebcca5-a963-40c5-a924-512f8562aeae" />







⚠️ Yasal Uyarı (Disclaimer)
Bu yazılım "Olduğu Gibi" (As Is) sunulmaktadır. Albion Online Kullanım Şartları (TOS), oyun trafiğinin dinlenmesi konusunda katı kurallara sahip olabilir. Bu yazılım herhangi bir oyun verisini değiştirmez (read-only), ancak kullanımı tamamen kullanıcının sorumluluğundadır.
