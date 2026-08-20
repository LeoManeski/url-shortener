import { useState, useEffect } from "react";

const API = "http://localhost:3000";

export default function App() {
  const [tab, setTab] = useState("shorten");
  const [url, setUrl] = useState("");
  const [desc, setDesc] = useState("");
  const [tags, setTags] = useState("");
  const [ttl, setTtl] = useState("");
  const [result, setResult] = useState(null);
  const [query, setQuery] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [recent, setRecent] = useState([]);
  const [top, setTop] = useState([]);
  const [statsCode, setStatsCode] = useState("");
  const [stats, setStats] = useState(null);

  const shorten = async () => {
    const body = { url };
    if (desc) body.description = desc;
    if (tags) body.tags = tags.split(",").map(t => t.trim()).filter(Boolean);
    if (ttl) body.ttl = parseInt(ttl);
    const res = await fetch(`${API}/shorten`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    setResult(await res.json());
    setUrl(""); setDesc(""); setTags(""); setTtl("");
  };

  const search = async () => {
    const res = await fetch(`${API}/search?q=${encodeURIComponent(query)}`);
    const data = await res.json();
    setSearchResults(data.results || []);
  };

  const fetchRecent = async () => {
    const res = await fetch(`${API}/recent`);
    const data = await res.json();
    setRecent(data.recent || []);
  };

  const fetchTop = async () => {
    const res = await fetch(`${API}/top`);
    const data = await res.json();
    setTop(data.top || []);
  };

  const fetchStats = async () => {
    const res = await fetch(`${API}/stats/${statsCode}`);
    const data = await res.json();
    setStats(data.error ? null : data);
  };

  useEffect(() => {
    if (tab === "recent") fetchRecent();
    if (tab === "top") fetchTop();
  }, [tab]);

  return (
    <div style={{ maxWidth: 550, margin: "30px auto", fontFamily: "system-ui", padding: 20 }}>
      <h2>URL Shortener</h2>

      <div style={{ display: "flex", gap: 8, marginBottom: 20 }}>
        {["shorten", "search", "recent", "top", "stats"].map(t => (
          <button key={t} onClick={() => setTab(t)}
            style={{ fontWeight: tab === t ? "bold" : "normal", cursor: "pointer" }}>
            {t}
          </button>
        ))}
      </div>

      {tab === "shorten" && (
        <div>
          <input value={url} onChange={e => setUrl(e.target.value)} placeholder="URL *" style={input} />
          <input value={desc} onChange={e => setDesc(e.target.value)} placeholder="Description (for search)" style={input} />
          <input value={tags} onChange={e => setTags(e.target.value)} placeholder="Tags (comma separated)" style={input} />
          <input value={ttl} onChange={e => setTtl(e.target.value)} placeholder="TTL in seconds" type="number" style={input} />
          <button onClick={shorten} style={btn}>Shorten</button>
          {result && (
            <div style={card}>
              <a href={result.short_url} target="_blank">{result.short_url}</a>
              <span style={{ marginLeft: 10, color: "#888", fontSize: 13 }}>({result.code})</span>
              <button onClick={() => navigator.clipboard.writeText(result.short_url)} style={{ marginLeft: 10, cursor: "pointer" }}>Copy</button>
            </div>
          )}
        </div>
      )}

      {tab === "search" && (
        <div>
          <input value={query} onChange={e => setQuery(e.target.value)} placeholder="Search by meaning..." style={input} onKeyDown={e => e.key === "Enter" && search()} />
          <button onClick={search} style={btn}>Search</button>
          {searchResults.map((r, i) => (
            <div key={i} style={card}>
              <div><a href={r.url} target="_blank">{r.url}</a></div>
              <div style={{ fontSize: 13, color: "#666" }}>{r.description}</div>
              {r.score && <div style={{ fontSize: 12, color: "#999" }}>Score: {parseFloat(r.score).toFixed(3)}</div>}
            </div>
          ))}
        </div>
      )}

      {tab === "recent" && (
        <div>
          <button onClick={fetchRecent} style={btn}>Refresh</button>
          {recent.map((code, i) => (
            <div key={i} style={card}>
              <a href={`${API}/${code}`} target="_blank">{API}/{code}</a>
              <button onClick={() => { setStatsCode(code); setTab("stats"); }} style={{ marginLeft: 10, cursor: "pointer" }}>Stats</button>
            </div>
          ))}
        </div>
      )}

      {tab === "top" && (
        <div>
          <button onClick={fetchTop} style={btn}>Refresh</button>
          {top.map((item, i) => (
            <div key={i} style={card}>
              <strong>#{i + 1}</strong> {API}/{item.Member} — <strong>{item.Score} clicks</strong>
            </div>
          ))}
        </div>
      )}

      {tab === "stats" && (
        <div>
          <input value={statsCode} onChange={e => setStatsCode(e.target.value)} placeholder="Link code" style={input} onKeyDown={e => e.key === "Enter" && fetchStats()} />
          <button onClick={fetchStats} style={btn}>Lookup</button>
          {stats && (
            <div style={card}>
              <div><strong>URL:</strong> {stats.url}</div>
              <div><strong>Clicks:</strong> {stats.clicks}</div>
              <div><strong>Unique visitors:</strong> {stats.unique_visitors}</div>
              <div><strong>Created:</strong> {stats.created}</div>
              <div><strong>TTL:</strong> {stats.ttl < 0 ? "No expiry" : `${Math.round(stats.ttl)}s`}</div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

const input = { width: "100%", padding: 8, marginBottom: 8, fontSize: 14, boxSizing: "border-box" };
const btn = { padding: "8px 16px", fontSize: 14, cursor: "pointer", marginBottom: 12 };
const card = { padding: 10, background: "#f5f5f5", borderRadius: 4, marginTop: 8 };