import React from 'react';

export const McpServersPage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4 flex justify-between items-center">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">Model Context Protocol Gateway</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Serveurs MCP & Outils Externes</h1>
        </div>
        <button className="bg-white hover:bg-zinc-200 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          + Connecter un Serveur MCP (Streamable-HTTP)
        </button>
      </div>

      <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4 text-xs">
        <div className="flex justify-between items-center border-b border-zinc-800/80 pb-3">
          <span className="font-semibold text-white">Serveurs MCP Connectes</span>
          <span className="text-xs font-mono text-cyan-400">Decouverte automatique tools/list</span>
        </div>

        <div className="space-y-3 font-mono text-[11px]">
          <div className="p-3 bg-zinc-900/60 rounded-xl border border-zinc-800 flex justify-between items-center">
            <div>
              <div className="font-bold text-white">Postgres DB MCP Server</div>
              <div className="text-zinc-400 text-[10px]">URL: https://mcp.internal/postgres/mcp | Outils: 4 decouverts</div>
            </div>
            <span className="px-2 py-0.5 rounded bg-emerald-950 text-emerald-300 border border-emerald-800 font-bold">CONNECTE</span>
          </div>
        </div>
      </div>
    </div>
  );
};
