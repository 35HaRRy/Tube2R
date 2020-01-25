using System.Collections.Generic;
using System.Net;
using System.Text.RegularExpressions;

using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore.Internal;

using Newtonsoft.Json;

namespace Tube2R.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class RSSController : ControllerBase
    {
        [HttpGet]
        public ActionResult Get([FromQuery] string channelId, [FromQuery] string playList)
        {
            var client = new WebClient();
            client.Headers.Add("Content-Type", "");

            var postValue = new Dictionary<string, object>()
            {
                { "url",  "https://www.youtube.com/" + (!string.IsNullOrEmpty(channelId) ? "channel/" + channelId : "playlist?list=" + playList) },
                { "format", "audio" },
                { "quality", "high" },
                { "page_size", 300 }
            };

            var newRSSResult = client.UploadString("http://podsync.net/api/create", JsonConvert.SerializeObject(postValue));
            var rssId = JsonConvert.DeserializeObject<Dictionary<string, string>>(newRSSResult);

            var rss = client.DownloadString("http://podsync.net/" + rssId["id"]);

            var format = "<guid>{0}</guid>";

            var guids = new Regex(@"<guid>(.*)<\/guid>").Matches(rss);
            var links = new Regex(@"<link>(.*)<\/link>").Matches(rss);

            foreach (Match guid in guids)
            {
                var lIndex = guids.IndexOf(guid) + 2;
                rss = rss.Replace(string.Format(format, guid.Groups[1].Value), string.Format(format, links[lIndex].Groups[1].Value));
            }

            return Content(rss, "text/xml");
        }
    }
}
