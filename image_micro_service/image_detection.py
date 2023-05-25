import json
import io
from keras.applications.resnet import (
    ResNet50,
    preprocess_input,
    decode_predictions,
)
from keras.preprocessing import image
from keras.utils import img_to_array
from PIL import Image
import numpy as np
import logging

log = logging.getLogger("ImageToKeywords")

# Image object detection model
model = ResNet50(weights="imagenet")


def image_to_keywords(file: str) -> tuple[str, int]:
    if file.filename == "":
        return ("No selected file", 400)
    log.info(file.filename)

    # Open image
    try:
        img = (
            Image.open(io.BytesIO(file.read()))
            .convert("RGB")
            .resize((224, 224))
        )
    except Exception as e:
        msg = "Failed to load image: " + str(e)
        log.error(msg)
        return (msg, 500)

    # Image to array
    try:
        x = img_to_array(img)
    except Exception as e:
        msg = "Failed to convert image to array: " + str(e)
        log.error(msg)
        return (msg, 500)

    # Image preprocess
    try:
        x = np.expand_dims(x, axis=0)
        x = preprocess_input(x)
    except Exception as e:
        msg = "Failed to preprocess image: " + str(e)
        log.error(msg)
        return (msg, 500)

    # Image detection
    try:
        preds = model.predict(x)
        predictions = decode_predictions(preds, top=5)[0]
        # Just return the top prediction
        keyword = predictions[0][1]
        log.info(keyword)
        return (json.dumps({"keyword": keyword}), 200)
    except Exception as e:
        msg = "Failed to process image: " + str(e)
        log.error(msg)
        return (msg, 500)
