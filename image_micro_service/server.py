from flask import Flask, request
from keras.applications.resnet import (
    ResNet50,
    preprocess_input,
    decode_predictions,
)
from keras.preprocessing import image
from keras.utils import img_to_array
from PIL import Image
import numpy as np
import io
import json
import logging

app = Flask(__name__)
model = ResNet50(weights="imagenet")

log = logging.getLogger("image_to_keywords")


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

    try:
        img = Image.open(io.BytesIO(file.read())).resize((224, 224))
    except Exception as e:
        msg = "Failed to open image: " + str(e)
        log.error(msg)
        return msg, 500

    try:
        x = img_to_array(img)
    except Exception as e:
        msg = "Failed to convert image to array: " + str(e)
        log.error(msg)
        return msg, 500

    try:
        x = np.expand_dims(x, axis=0)
        x = preprocess_input(x)
    except Exception as e:
        msg = "Failed to preprocess image: " + str(e)
        log.error(msg)
        return msg, 500

    try:
        preds = model.predict(x)
        predictions = decode_predictions(preds, top=5)[0]
        # Just return the top prediction
        keyword = predictions[0][1]
        log.info(keyword)

        return json.dumps({"keyword": keyword}), 200
    except Exception as e:
        msg = "Failed to process image: " + str(e)
        log.error(msg)
        return msg, 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9021)
