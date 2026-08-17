import React, { useState } from 'react';

export const LandingPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'voice' | 'whatsapp' | 'visio' | 'mcp'>('voice');

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 font-sans p-6 space-y-12 max-w-7xl mx-auto">
      {/* Header / Hero */}
      <div className="text-center space-y-4 pt-12">
        <span className="text-xs font-mono font-bold text-cyan-400 bg-cyan-950/80 px-3 py-1.5 rounded-full border border-cyan-800/80">
          Plateforme Conversationnelle IA Sovereign & Multi-Canal
        </span>
        <h1 className="text-4xl sm:text-6xl font-extrabold tracking-tight text-zinc-100">
          KallFlow <span className="text-cyan-400">Conversational Engine</span>
        </h1>
        <p className="text-sm sm:text-base text-zinc-400 max-w-2xl mx-auto leading-relaxed">
          Infrastructures temps réel sous sub-500ms pour Téléphonie Vocale, WhatsApp, Visio Bot (Zoom/Teams/Meet) et Gateway MCP.
        </p>
      </div>

      {/* Simulator Preview Tabs */}
      <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-2xl space-y-6">
        <div className="flex border-b border-zinc-800 pb-3 gap-3 text-xs font-bold overflow-x-auto">
          <button
            onClick={() => setActiveTab('voice')}
            className={`px-4 py-2 rounded-xl transition-all ${
              activeTab === 'voice' ? 'bg-cyan-500 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white bg-zinc-900'
            }`}
          >
            Agent Vocal Telephonique
          </button>
          <button
            onClick={() => setActiveTab('whatsapp')}
            className={`px-4 py-2 rounded-xl transition-all ${
              activeTab === 'whatsapp' ? 'bg-cyan-500 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white bg-zinc-900'
            }`}
          >
            Messaging WhatsApp & RCS
          </button>
          <button
            onClick={() => setActiveTab('visio')}
            className={`px-4 py-2 rounded-xl transition-all ${
              activeTab === 'visio' ? 'bg-cyan-500 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white bg-zinc-900'
            }`}
          >
            Bot Visio & Copilote
          </button>
          <button
            onClick={() => setActiveTab('mcp')}
            className={`px-4 py-2 rounded-xl transition-all ${
              activeTab === 'mcp' ? 'bg-cyan-500 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white bg-zinc-900'
            }`}
          >
            Gateway MCP & Outillage
          </button>
        </div>

        {/* Tab Content */}
        <div className="p-4 bg-zinc-900/60 rounded-xl border border-zinc-800 text-xs font-mono space-y-3">
          {activeTab === 'voice' && (
            <div className="space-y-2">
              <div className="text-cyan-400 font-bold">▶ Test Simulateur Vocal Pipeline Sub-350ms</div>
              <p className="text-zinc-300">STT Neural Ultra-Fast + LLM Router + High-Fidelity TTS avec annulation de barge-in G711/PCM.</p>
              <div className="p-3 bg-zinc-950 rounded border border-zinc-800 text-zinc-400">
                [Agent Vocal]: "Bonjour ! Je suis l'assistant d'accueil KallFlow. Comment puis-je vous guider ?"
              </div>
            </div>
          )}

          {activeTab === 'whatsapp' && (
            <div className="space-y-2">
              <div className="text-cyan-400 font-bold">▶ Validation JID WhatsApp & Fallback SMS/RCS</div>
              <p className="text-zinc-300">Normalisation E.164 et dispatch automatique avec jitter anti-spam (+/- 30%).</p>
              <div className="p-3 bg-zinc-950 rounded border border-zinc-800 text-emerald-400">
                [Status]: Contact validé E.164 (+33612345678@s.whatsapp.net) | Session Active.
              </div>
            </div>
          )}

          {activeTab === 'visio' && (
            <div className="space-y-2">
              <div className="text-cyan-400 font-bold">▶ Bot Visio Multi-Plateformes Native</div>
              <p className="text-zinc-300">Connexion Zoom RTMS, Google Meet et Microsoft Teams avec isolation audio multi-pistes.</p>
              <div className="p-3 bg-zinc-950 rounded border border-zinc-800 text-indigo-300">
                [Meeting Engine]: Bot identifié connecté | Enregistrement HD & Diarisation en direct.
              </div>
            </div>
          )}

          {activeTab === 'mcp' && (
            <div className="space-y-2">
              <div className="text-cyan-400 font-bold">▶ Protocol Model Context Protocol (MCP Streamable-HTTP)</div>
              <p className="text-zinc-300">Découverte automatique `tools/list` pour HubSpot, Salesforce, Postgres et outils métiers sur-mesure.</p>
              <div className="p-3 bg-zinc-950 rounded border border-zinc-800 text-amber-300">
                [MCP Gateway]: 10 outils découverts et isolés par TenantID RLS.
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Feature Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 space-y-3">
          <h3 className="font-bold text-white text-base">Moteur Vocal Native</h3>
          <p className="text-xs text-zinc-400">Latence ultra-faible avec gestion du barge-in et suppression du bruit Krisp/RNNoise.</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 space-y-3">
          <h3 className="font-bold text-white text-base">Bot Visio Native</h3>
          <p className="text-xs text-zinc-400">Capture vidéo HD 1080p et isolation PCM separee pour diarisation parfaite.</p>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 space-y-3">
          <h3 className="font-bold text-white text-base">Multi-Tenancy RLS & MCP</h3>
          <p className="text-xs text-zinc-400">Isolation PostgreSQL Row-Level Security et outils MCP isolés par tenant.</p>
        </div>
      </div>
    </div>
  );
};
