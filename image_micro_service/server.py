from flask import Flask, request
import json
import logging
import entity_detection
import image_detection


log = logging.getLogger("image_to_keywords")

app = Flask(__name__)


@app.route("/image_to_keywords", methods=["POST"])
def image_to_keywords():
    # Check if the post request has the file part
    if "file" not in request.files:
        return "No file part", 400
    file = request.files["file"]

    # if user does not select file, browser submits an empty part without
    # filename
    if file.filename == "":
        return "No selected file", 400
    log.info(file.filename)
    result, code = image_detection.image_to_keywords(file)
    if code != 200:
        log.error(result)
        return result, code
    return result, code


@app.route("/extract_info", methods=["POST"])
def extract_info():
    data = request.get_json()
    text = data.get("text")
    language = data.get("language")

    if not text or not language:
        return "Invalid request: no text or no language", 400

    if language not in ["en", "cn"]:
        return "Unsupported language: " + language, 400

    # Entity detection
    entities = entity_detection.entity_detect(text, language)

    return json.dumps({"entities": entities})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9021)
