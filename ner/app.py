
from flask import Flask, jsonify, request
from flask_cors import CORS
import spacy

app = Flask(__name__)
CORS(app)

nlp = spacy.load("en_core_web_sm")

@app.route("/extract-entities", methods=["POST"])
def extract_entities():
    try:
        req_data = request.get_json()

        if isinstance(req_data, dict) and 'prompt' in req_data:
            text = req_data['prompt']
        else:
            return jsonify({"status": "failure", "message": "Invalid input format"}), 400


        if text:
            text = text.strip()
            doc = nlp(text)
            locations = []
            from_location = None
            to_location = None

            for ent in doc.ents:
                if ent.label_ == "GPE":
                    locations.append(ent.text)

            for token in doc:
                if token.dep_ == "pobj" and token.head.text in ["from", "to"]:
                    if token.head.text == "from":
                        from_location = token.text
                    else:
                        to_location = token.text

            result = {
                "locations": locations,
                "from": from_location,
                "to": to_location
            }

            return jsonify(result), 200
        else:
            return jsonify({"status": "failure", "message": "no input provided"}), 400
    except Exception as e:
        return jsonify({"status": "failure", "message": str(e)}), 500

if __name__ == "__main__":
    app.run(port=2525)
