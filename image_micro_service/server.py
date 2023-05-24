from flask import Flask, request
from keras.applications.resnet import (
    ResNet50,
    preprocess_input,
    decode_predictions,
)
from keras.preprocessing import image
from PIL import Image
import numpy as np
import io
import json

app = Flask(__name__)
model = ResNet50(weights="imagenet")


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

    try:
        img = Image.open(io.BytesIO(file.read())).resize((224, 224))
        x = image.img_to_array(img)
        x = np.expand_dims(x, axis=0)
        x = preprocess_input(x)

        preds = model.predict(x)
        predictions = decode_predictions(preds, top=5)[0]

        # Just return the top prediction
        keyword = predictions[0][1]

        return json.dumps({"keyword": keyword}), 200

    except Exception as e:
        return "Failed to process image", 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=9021)
