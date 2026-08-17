import React, { useState } from 'react';

export const LandingPage: React.FC = () => {
  const [simulatingCall, setSimulatingCall] = useState(false);
  const [activeTab, setActiveTab] = useState<'voice' | 'whatsapp' | 'visio' | 'mcp' | 'avatar'>('voice');

  const handleSimulateCall = () => {
    setSimulatingCall(true);
    setTimeout(() => setSimulatingCall(false), 4000);
  };

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 font-sans selection:bg-cyan-500 selection:text-zinc-950">
      {/* Top Showcase Navigation */}
      <header className="border-b border-zinc-800/80 sticky top-0 z-50 bg-[#09090b]/80 backdrop-blur-md">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between text-xs">
          <div className="flex items-center gap-3">
            <span className="text-xl font-black bg-clip-text text-transparent bg-gradient-to-r from-violet-400 to-cyan-400">
              KallFlow
            </span>
            <span className="text-[10px] bg-zinc-900 border border-zinc-800 text-zinc-400 px-2 py-0.5 rounded font-mono font-bold">
              v0.7.1 Enterprise
            </span>
          </div>

          <nav className="hidden md:flex items-center gap-6 text-zinc-400 font-medium">
            <a href="#features" className="hover:text-zinc-100 transition-colors">Plateforme & Pipeline</a>
            <a href="#campaigns" className="hover:text-zinc-100 transition-colors">Campagnes Omnicanales</a>
            <a href="#mcp" className="hover:text-zinc-100 transition-colors">Passerelle MCP & RAG</a>
            <a href="#pricing" className="hover:text-zinc-100 transition-colors">Tarifs & Credits</a>
          </nav>

          <div className="flex items-center gap-3">
            <button className="bg-white hover:bg-zinc-200 text-zinc-950 font-semibold px-4 py-2 rounded-lg transition-colors shadow-sm">
              Lancer la Console →
            </button>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="max-w-7xl mx-auto px-6 pt-20 pb-16 text-center space-y-6">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-zinc-900 border border-zinc-800 text-xs text-zinc-300 font-mono">
          <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          KallFlow Engine: Voice, Visio, WhatsApp & Avatars GPU Local
        </div>

        <h1 className="text-4xl sm:text-6xl font-bold tracking-tight text-zinc-100 max-w-4xl mx-auto leading-tight">
          La Infrastructure Conversationnelle Unifiee pour Agents IA Enterprise
        </h1>

        <p className="text-sm sm:text-base text-zinc-400 max-w-2xl mx-auto leading-relaxed">
          Pipelining temps reel (STT Deepgram + LLM OpenRouter + TTS ElevenLabs), VAD G711/PCM, annulation de barge-in thread-safe, bots visio Recall.ai, campagnes WhatsApp & SMS anti-spam et passerelle MCP Postgres RLS.
        </p>

        <div className="pt-4 flex flex-col sm:flex-row justify-center items-center gap-4 text-xs font-semibold">
          <button className="w-full sm:w-auto bg-white hover:bg-zinc-200 text-zinc-950 px-6 py-3 rounded-xl transition-colors shadow-lg">
            Demarrer un Essai Gratuit
          </button>
          <a href="#demo" className="w-full sm:w-auto bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-200 px-6 py-3 rounded-xl transition-colors">
            Tester la Demo Live
          </a>
        </div>
      </section>

      {/* Interactive Live Simulator Demo */}
      <section id="demo" className="max-w-5xl mx-auto px-6 py-12">
        <div className="bg-zinc-900/60 border border-zinc-800 rounded-2xl p-6 shadow-2xl backdrop-blur-md space-y-6">
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center border-b border-zinc-800/80 pb-4 gap-4">
            <div>
              <h3 className="text-base font-semibold text-zinc-100">Simulateur Moteur Temps Reel</h3>
              <p className="text-xs text-zinc-400 mt-0.5">Testez la latence sub-500ms, la VAD PCM et le barge-in barge-in instantane.</p>
            </div>

            <div className="flex bg-zinc-950 p-1 border border-zinc-800 rounded-xl text-xs">
              {(['voice', 'whatsapp', 'visio', 'mcp', 'avatar'] as const).map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`px-3 py-1.5 rounded-lg transition-colors font-medium capitalize ${
                    activeTab === tab ? 'bg-zinc-100 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white'
                  }`}
                >
                  {tab}
                </button>
              ))}
            </div>
          </div>

          {/* Tab Content Display */}
          <div className="bg-zinc-950 p-6 rounded-xl border border-zinc-800 space-y-4 font-mono text-xs">
            {activeTab === 'voice' && (
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-zinc-400">[Pipeline Mode] G711-to-PCM VAD + Smart-Turn v3 + Barge-In Cancel</span>
                  <span className="text-emerald-400">TTFA: 320ms | Cost Opti Matrix Active</span>
                </div>
                <div className="p-3 bg-zinc-900/60 rounded border border-zinc-800/80 text-zinc-200">
                  <span className="text-cyan-400 font-bold">Agent KallFlow:</span> "Bonjour, je suis l'assistant vocal KallFlow Engine. Comment puis-je vous aider aujourd'hui ?"
                </div>
                <button
                  onClick={handleSimulateCall}
                  disabled={simulatingCall}
                  className="bg-cyan-500 hover:bg-cyan-400 text-zinc-950 font-bold px-4 py-2 rounded-lg transition-colors text-xs"
                >
                  {simulatingCall ? 'Session Audio Active...' : 'Lancer une Session Vocale (Simulee)'}
                </button>
              </div>
            )}

            {activeTab === 'whatsapp' && (
              <div className="space-y-3">
                <div className="p-3 bg-zinc-900/60 rounded border border-zinc-800/80 text-zinc-200">
                  <span className="text-emerald-400 font-bold">WhatsApp Engine:</span> Normalisation E.164 (<code className="text-cyan-400">@s.whatsapp.net</code>), rotation de session round-robin et fallback SMS automatique.
                </div>
              </div>
            )}

            {activeTab === 'visio' && (
              <div className="space-y-3">
                <div className="p-3 bg-zinc-900/60 rounded border border-zinc-800/80 text-zinc-200">
                  <span className="text-violet-400 font-bold">Meeting Engine:</span> Integration Recall.ai pour Zoom RTMS, Google Meet Media API & Teams Signed-in Bots avec scoring BANT LeMUR.
                </div>
              </div>
            )}

            {activeTab === 'mcp' && (
              <div className="space-y-3">
                <div className="p-3 bg-zinc-900/60 rounded border border-zinc-800/80 text-zinc-200">
                  <span className="text-amber-400 font-bold">Passerelle MCP:</span> Decouverte dynamique <code className="text-cyan-400">tools/list</code> et execution sous PostgreSQL RLS avec isolation tenant.
                </div>
              </div>
            )}

            {activeTab === 'avatar' && (
              <div className="space-y-3">
                <div className="p-3 bg-zinc-900/60 rounded border border-zinc-800/80 text-zinc-200">
                  <span className="text-rose-400 font-bold">Avatar Engine:</span> Mode hybride WaveSpeed API temps reel & conteneurs <code className="text-cyan-400">avatar-renderer-mcp</code> sur pool GPU local.
                </div>
              </div>
            )}
          </div>
        </div>
      </section>

      {/* Feature Pillars */}
      <section id="features" className="max-w-7xl mx-auto px-6 py-16 border-t border-zinc-800/80">
        <h2 className="text-2xl font-bold tracking-tight text-zinc-100 text-center mb-12">
          Architecture Moteur Conversationnel KallFlow
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-xs">
          <div className="bg-zinc-900/40 border border-zinc-800/80 rounded-xl p-5 space-y-2">
            <h4 className="font-bold text-zinc-100 text-sm">Gouvernance Multi-Tenant 3-Tier</h4>
            <p className="text-zinc-400 leading-relaxed">
              Super Admin Platform, Agence Marque Blanche (Plan Agency avec marge personnalisee) et Tenants Enterprise avec gestion des roles stricte.
            </p>
          </div>

          <div className="bg-zinc-900/40 border border-zinc-800/80 rounded-xl p-5 space-y-2">
            <h4 className="font-bold text-zinc-100 text-sm">Sélecteur d'Intégration Hybride</h4>
            <p className="text-zinc-400 leading-relaxed">
              Connexion 1-Click KallFlow Managed ou propre application developpeur Meta/TikTok (BYO-App) pour un anonymat total.
            </p>
          </div>

          <div className="bg-zinc-900/40 border border-zinc-800/80 rounded-xl p-5 space-y-2">
            <h4 className="font-bold text-zinc-100 text-sm">Systeme de Credits Unifie</h4>
            <p className="text-zinc-400 leading-relaxed">
              Equivalence stricte 1 EUR = 1 000 Credits KallFlow avec debit en temps reel sur les conteneurs vocal, visio et messaging.
            </p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-zinc-800/80 py-8 text-center text-xs text-zinc-500">
        <p>© 2026 KallFlow Inc. Tous droits reserves. (https://kallflow.ai)</p>
      </footer>
    </div>
  );
};
