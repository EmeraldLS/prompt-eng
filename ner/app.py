
from flask import Flask, jsonify, request
import spacy

app = Flask(__name__)
nlp = spacy.load("en_core_web_sm")

@app.route("/extract-entities", methods=["POST"])
def extract_entities():
    try:
        req_data = request.get_json()
        text = req_data.get("text")

        if text:
            doc = nlp(text)
            locations = [ ent.text for ent in doc.ents if ent.label_ == "GPE"]
            
            for token in doc:
                
                print(token.text, token.pos_, token.lemma_)

            return jsonify(locations), 200
        else:
            return jsonify({"status": "failure", "message": "no input provided"}), 400
    except Exception as e:
        return jsonify({"status": "failure", "message": str(e)}), 500

if __name__ == "__main__":
    app.run(port=2332)
