import datetime
import scrapy
import logging

logger = logging.getLogger("OiWikiSpider")
logger.setLevel(logging.DEBUG)


class OiwikiSpiderSpider(scrapy.Spider):
    name = "oiwiki_spider"
    allowed_domains = ["oi-wiki.org"]
    start_urls = ["https://oi-wiki.org/"]
    local_site = "oi-wiki_response.html"
    id = 0

    def start_requests(self):
        urls = [
            "https://oi-wiki.org/contest/",
            "https://oi-wiki.org/tools/",
            "https://oi-wiki.org/lang/",
            "https://oi-wiki.org/basic/",
            "https://oi-wiki.org/search/",
            "https://oi-wiki.org/dp/",
            "https://oi-wiki.org/string/",
            "https://oi-wiki.org/math/",
            "https://oi-wiki.org/ds/",
            "https://oi-wiki.org/graph/",
            "https://oi-wiki.org/geometry/",
            "https://oi-wiki.org/misc/",
            "https://oi-wiki.org/misc/",
        ]

        for i, url in enumerate(urls):
            yield scrapy.Request(url=url, callback=self.parse, cb_kwargs={})

    def parse(self, response):
        # If there's two-level selector, yield scrapy.Request separately to its
        # leaf nodes
        sections = response.xpath(
            "//li[@class='md-nav__item']/a[@class='md-nav__link']"
        )
        hrefs = sections.xpath("@href").getall()
        texts = sections.xpath("text()").getall()
        texts = [t.strip() for t in texts]
        # Number of hrefs should be equal to number of texts
        assert len(hrefs) == len(texts)

        # Yield scrapy request to all sections
        for href, section in zip(hrefs, texts):
            url = response.urljoin(href)
            yield scrapy.Request(
                url=url,
                callback=self.parse_section,
                cb_kwargs={"section": section},
            )

    def parse_section(self, response, section="Unknown"):
        # Scrapy section details
        content = response.xpath(
            '//div[@class="md-content"]//*[self::h1 or self::h2 or self::h3 or self::h4 or self::li or self::ul or self::p]/text()'
        ).getall()

        self.id = self.id+1
        yield {
            "id": self.id,
            "title": content[0],
            "content": "\n".join(para for para in content),
            "url": response.url,
            "date": datetime.date.today().strftime("%Y-%m-%d"),
        }
