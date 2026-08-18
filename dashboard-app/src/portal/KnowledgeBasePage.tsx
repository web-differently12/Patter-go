import React from 'react';

export const KnowledgeBasePage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4 flex justify-between items-center">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">RAG Knowledge Base & pgvector Store</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Base de Connaissances & Vector Store</h1>
        </div>
        <button className="bg-white hover:bg-zinc-200 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          + Ingerer un Document (PDF/TXT)
        </button>
      </div>

      <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4 text-xs">
        <div className="flex justify-between items-center border-b border-zinc-800/80 pb-3">
          <span className="font-semibold text-white">Documents Ingeres (PostgreSQL pgvector)</span>
          <span className="text-xs font-mono text-cyan-400">Embeddings: text-embedding-3-small (1536 dim)</span>
        </div>

        <div className="space-y-3">
          <div className="p-3 bg-zinc-900/60 rounded-xl border border-zinc-800 flex justify-between items-center">
            <div>
              <div className="font-bold text-zinc-200">documentation_api_v2.pdf</div>
              <div className="text-[11px] text-zinc-500 font-mono">142 chunks vectorises · Index Cosine</div>
            </div>
            <span className="text-xs bg-emerald-950 text-emerald-300 border border-emerald-800 px-2 py-0.5 rounded font-mono font-bold">ACTIF</span>
          </div>
        </div>
      </div>
    </div>
  );
};
