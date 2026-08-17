import React, { useState } from 'react';

export const AgentEditorPage: React.FC = () => {
  const [systemPrompt, setSystemPrompt] = useState(
    'Vous etes l\'assistant accueil d\'entreprise. Repondez de maniere courtoise, synthetique et sans jargon.'
  );
  const [model, setModel] = useState('Neural Realtime Engine');
  const [voice, setVoice] = useState('Voice Flash Ultra-HD');
  const [bargeInThresholdMs, setBargeInMs] = useState(300);

  // Vapi-inspired Voice Agent Tuning
  const [backchanneling, setBackchanneling] = useState(true);
  const [ambientNoiseSound, setAmbientNoise] = useState<'office' | 'cafe' | 'silent'>('office');
  const [fillerWords, setFillerWords] = useState(true);

  // LockedInAI-inspired Realtime Visio Copilot Tuning
  const [copilotSuggestions, setCopilotSuggestions] = useState(true);
  const [objectionDetector, setObjectionDetector] = useState(true);

  // Lindy.ai-inspired Retail Workflow Triggers
  const [crmAutoSync, setCrmAutoSync] = useState(true);
  const [calendarAutoSchedule, setCalendarAutoSchedule] = useState(true);

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
          <label className="block text-xs font-semibold text-zinc-300">System Prompt et Instructions Telephoniques</label>
          <textarea
            value={systemPrompt}
            onChange={(e) => setSystemPrompt(e.target.value)}
            rows={10}
            className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-xs text-zinc-200 font-mono focus:outline-none focus:border-cyan-500 leading-relaxed"
          />

          {/* LockedInAI & Lindy Workflow Options */}
          <div className="border-t border-zinc-800/80 pt-4 space-y-3 text-xs">
            <h4 className="font-bold text-white">Copilote Temps Reel & Workflows Automatises</h4>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 font-mono">
              <label className="flex items-center gap-2 p-2 bg-zinc-900 rounded border border-zinc-800">
                <input
                  type="checkbox"
                  checked={copilotSuggestions}
                  onChange={(e) => setCopilotSuggestions(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Suggestions Copilote en direct</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-zinc-900 rounded border border-zinc-800">
                <input
                  type="checkbox"
                  checked={objectionDetector}
                  onChange={(e) => setObjectionDetector(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Detecteur d'objections B2B</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-zinc-900 rounded border border-zinc-800">
                <input
                  type="checkbox"
                  checked={crmAutoSync}
                  onChange={(e) => setCrmAutoSync(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Auto-Sync CRM / Fiche Client</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-zinc-900 rounded border border-zinc-800">
                <input
                  type="checkbox"
                  checked={calendarAutoSchedule}
                  onChange={(e) => setCalendarAutoSchedule(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Prise de RDV Calendrier auto</span>
              </label>
            </div>
          </div>
        </div>

        {/* Engine, Voice & Tuning */}
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4 text-xs">
          <h3 className="font-bold text-white text-sm">Paramètres Moteur et Voix</h3>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Modele LLM / Realtime</label>
            <select
              value={model}
              onChange={(e) => setModel(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Neural Realtime Engine">Neural Realtime Engine</option>
              <option value="Enterprise Intelligence Tier">Enterprise Intelligence Tier</option>
              <option value="Realtime Sales Copilot Engine">Realtime Sales Copilot Engine</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Voix Synthetique TTS</label>
            <select
              value={voice}
              onChange={(e) => setVoice(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Voice Flash Ultra-HD">Voice Flash Ultra-HD</option>
              <option value="Voice Sonic HD">Voice Sonic HD</option>
              <option value="Voice Alloy HD">Voice Alloy HD</option>
            </select>
          </div>

          {/* Vapi-inspired Fine Voice Options */}
          <div className="border-t border-zinc-800/80 pt-3 space-y-2 font-mono">
            <span className="font-bold text-zinc-300 block">Fine Conversational Options</span>

            <label className="flex items-center gap-2 text-[11px]">
              <input
                type="checkbox"
                checked={backchanneling}
                onChange={(e) => setBackchanneling(e.target.checked)}
                className="accent-cyan-500"
              />
              <span>Backchanneling ("mm-hmm", "d'accord")</span>
            </label>

            <label className="flex items-center gap-2 text-[11px]">
              <input
                type="checkbox"
                checked={fillerWords}
                onChange={(e) => setFillerWords(e.target.checked)}
                className="accent-cyan-500"
              />
              <span>Mots de remplissage pendant RAG</span>
            </label>

            <div>
              <label className="block text-zinc-400 text-[11px] mb-1">Bruit d'ambiance d'arriere-plan</label>
              <select
                value={ambientNoiseSound}
                onChange={(e) => setAmbientNoise(e.target.value as any)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1 text-zinc-300"
              >
                <option value="office">Ambiance Bureau Discret</option>
                <option value="cafe">Ambiance Cafe</option>
                <option value="silent">Silencieux (Insonorise)</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Seuil Interruption Barge-In: {bargeInThresholdMs}ms</label>
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
