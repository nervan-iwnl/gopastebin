// frontend/src/pages/PasteView.jsx
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { api } from "../api";

export default function PasteView() {
  const { slug } = useParams();
  const [paste, setPaste] = useState(null);
  const [raw, setRaw] = useState("");
  const [err, setErr] = useState("");

  useEffect(() => {
    async function load() {
      try {
        const data = await api.getPaste(slug);
        setPaste(data.paste || data);
        const txt = await api.getPasteRaw(slug);
        setRaw(txt);
      } catch (e) {
        setErr("not found");
      }
    }
    load();
  }, [slug]);

  if (err) return <p>{err}</p>;
  return (
    <div>
      <div className="card">
        <div className="card-body">
          <h4 className="card-title">
            {paste?.title} <small className="text-muted">{slug}</small>
          </h4>
          <div className="mb-3">
            <pre className="p-3" style={{ background: "#f8f9fa", whiteSpace: "pre-wrap" }}>
              {raw}
            </pre>
          </div>
          <div className="d-flex gap-2">
            <a className="btn btn-outline-secondary btn-sm" href={`/api/v1/pastes/${slug}/raw`} target="_blank" rel="noreferrer">Open raw</a>
            <button
              className="btn btn-secondary btn-sm"
              onClick={() => navigator.clipboard?.writeText(raw)}
            >
              Copy
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
