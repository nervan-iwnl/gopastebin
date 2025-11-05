// frontend/src/pages/MyPastes.jsx
import { useEffect, useState } from "react";
import { api } from "../api";
import PasteForm from "../components/PasteForm.jsx";
import { useAuth } from "../hooks/useAuth.jsx";

export default function MyPastes() {
  const [pastes, setPastes] = useState([]);
  const [error, setError] = useState("");

  async function load() {
    try {
      let data;
      if (user) {
        data = await api.myPastes();
        setPastes(data.pastes || data || []);
      } else {
        // guest: show recent public pastes
        data = await api.getRecent();
        setPastes(data.pastes || data || []);
      }
      setError("");
    } catch (e) {
      setError(e?.error?.message || "failed to load");
    }
  }

  const { user } = useAuth();

  useEffect(() => {
    load();
  }, [user]);

  return (
    <div>
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h2 className="m-0">My pastes</h2>
      </div>
      {error && <div className="alert alert-danger">{error}</div>}

      <div className="row">
        <div className="col-md-5">
          <PasteForm onCreated={load} />
          <div className="mt-3">
            <h6>Quick anonymous paste</h6>
            <div className="d-flex gap-2">
              <button
                className="btn btn-outline-secondary"
                onClick={async () => {
                    try {
                      const res = await api.createAnon("print('anon from ui')");
                      // response should include paste or slug
                      const slug = res?.paste?.slug || res?.slug || (res?.pastes && res.pastes[0]?.slug);
                      if (slug) {
                        // open paste view in new tab
                        window.open(`/p/${slug}`, "_blank");
                        // also open raw in new tab optionally
                        // window.open(`/api/v1/pastes/${slug}/raw`, "_blank");
                      } else {
                        alert('Anonymous paste created');
                      }
                      load();
                    } catch (e) {
                      alert(e?.error?.message || 'failed');
                    }
                  }}
              >
                Create anonymous paste
              </button>
            </div>
          </div>
        </div>
        <div className="col-md-7">
          <div className="list-group">
            {pastes.length === 0 && <div className="text-muted">No pastes yet</div>}
            {pastes.map((p) => (
              <a key={p.slug} href={`/p/${p.slug}`} className="list-group-item list-group-item-action">
                <div className="d-flex w-100 justify-content-between">
                  <h5 className="mb-1">{p.title || p.slug}</h5>
                  <small className="text-muted">{p.folder}</small>
                </div>
                <p className="mb-1 text-truncate">{p.summary || p.content || ''}</p>
                <small className="text-muted">{p.slug}</small>
              </a>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
