import React, { useState } from 'react';

export const AgentEditorPage: React.FC = () => {
  const [systemPrompt, setSystemPrompt] = useState(
    'Vous êtes l\'assistant accueil d\'entreprise. Répondez de manière courtoise, synthétique et sans jargon.'
  );
  const [model, setModel] = useState('gpt-realtime-mini');
  const [voice, setVoice] = useState('alloy');
  const [bargeInThresholdMs, setBargeInMs] = useState(300);

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4 flex justify-between items-center">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">Agent Studio Configurator</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Configuration Avancee de l'Agent</h1>
        </div>
        <button className="bg-cyan-500 hover:bg-cyan-400 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          Enregistrer l'Agent
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* System Prompt & Instructions */}
        <div className="lg:col-span-2 bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4">
          <label className="block text-xs font-semibold text-zinc-300">System Prompt & Instructions Telephoniques</label>
          <textarea
            value={systemPrompt}
            onChange={(e) => setSystemPrompt(e.target.value)}
            rows={12}
            className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-xs text-zinc-200 font-mono focus:outline-none focus:border-cyan-500 leading-relaxed"
          />
        </div>

        {/* Engine, Voice & Tuning */}
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4 text-xs">
          <h3 className="font-bold text-white text-sm">Paramètres Moteur & Voix</h3>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Modèle LLM / Realtime</label>
            <select
              value={model}
              onChange={(e) => setModel(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="gpt-realtime-mini">OpenAI Realtime Mini</option>
              <option value="claude-3-5-sonnet">Anthropic Claude 3.5 Sonnet</option>
              <option value="grok-voice">xAI Grok Voice Agent</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Voix Synthetique TTS</label>
            <select
              value={voice}
              onChange={(e) => setVoice(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="alloy">OpenAI Alloy</option>
              <option value="eleven_flash_v2_5">ElevenLabs Flash v2.5</option>
              <option value="cartesia_sonic">Cartesia Sonic</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Seuil Interruption Barge-In (ms): {bargeInThresholdMs}ms</label>
            <input
              type="range"
              min={100}
              max={1000}
              step={50}
              value={bargeInThresholdMs}
              onChange={(e) => setBargeInMs(Number(e.target.value))}
              className="w-full h-1.5 bg-zinc-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
            />
          </div>
        </div>
      </div>
    </div>
  );
};
