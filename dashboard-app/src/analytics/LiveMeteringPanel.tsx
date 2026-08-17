import React, { useState } from 'react';

export interface LiveTelemetryData {
  llmTokensIn: number;
  llmTokensOut: number;
  sttDurationMs: number;
  ttsCharacters: number;
  sipDurationMs: number;
  activeCallsCount: number;
}

export const LiveMeteringPanel: React.FC = () => {
  const [telemetry, setTelemetry] = useState<LiveTelemetryData>({
    llmTokensIn: 1420,
    llmTokensOut: 280,
    sttDurationMs: 42000, // 42.0s
    ttsCharacters: 8900,
    sipDurationMs: 125000, // 125.0s
    activeCallsCount: 4,
  });

  const [simulating, setSimulating] = useState(false);

  const handleSimulateTraffic = () => {
    setSimulating(true);
    setTelemetry((prev) => ({
      llmTokensIn: prev.llmTokensIn + Math.floor(Math.random() * 150),
      llmTokensOut: prev.llmTokensOut + Math.floor(Math.random() * 40),
      sttDurationMs: prev.sttDurationMs + 3500,
      ttsCharacters: prev.ttsCharacters + 250,
      sipDurationMs: prev.sipDurationMs + 5000,
      activeCallsCount: Math.max(1, prev.activeCallsCount + (Math.random() > 0.5 ? 1 : -1)),
    }));
    setTimeout(() => setSimulating(false), 300);
  };

  const formattedSttSeconds = (telemetry.sttDurationMs / 1000).toFixed(1);
  const formattedSipSeconds = (telemetry.sipDurationMs / 1000).toFixed(1);
  const totalTokens = telemetry.llmTokensIn + telemetry.llmTokensOut;
  const formattedTotalTokens = totalTokens >= 1000 ? `${(totalTokens / 1000).toFixed(1)}k` : `${totalTokens}`;

  return (
    <div className="bg-zinc-900/50 backdrop-blur-md border border-zinc-800/80 rounded-xl p-5 hover:border-zinc-700/80 transition-all duration-150 space-y-4 font-sans text-xs">
      <div className="flex justify-between items-center border-b border-zinc-800/60 pb-3">
        <div>
          <div className="flex items-center gap-2">
            <h3 className="text-sm font-semibold tracking-tight text-zinc-100">Live Metering - KallFlow Engine</h3>
            <span className="w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse"></span>
          </div>
          <p className="text-[11px] text-zinc-400 mt-0.5">Telemetrie en direct des micro-transactions LLM, STT, TTS et Telephonie SIP.</p>
        </div>

        {/* Outline Simulation Debug Button */}
        <button
          onClick={handleSimulateTraffic}
          disabled={simulating}
          className="border border-zinc-700/80 bg-zinc-800/40 hover:bg-zinc-800 text-zinc-300 text-xs font-medium px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1.5 disabled:opacity-50"
        >
          <span>{simulating ? 'Envoi...' : 'Simuler Trafic'}</span>
        </button>
      </div>

      {/* 4-Column Live Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Column 1: LLM Tokens */}
        <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-3.5 space-y-1.5">
          <div className="flex justify-between items-center text-zinc-400 font-medium">
            <span>LLM Brain Tier</span>
            <span className="text-[10px] bg-zinc-800 text-zinc-300 px-1.5 py-0.5 rounded font-mono">OpenRouter</span>
          </div>
          <div className="text-base font-semibold text-zinc-100 font-mono">
            {telemetry.llmTokensIn.toLocaleString()} In / {telemetry.llmTokensOut.toLocaleString()} Out
          </div>
          <div className="text-[11px] text-zinc-500 font-mono">
            ({formattedTotalTokens} total tokens)
          </div>
        </div>

        {/* Column 2: STT Speech-to-Text */}
        <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-3.5 space-y-1.5">
          <div className="flex justify-between items-center text-zinc-400 font-medium">
            <span>STT Transcription</span>
            <span className="text-[10px] bg-zinc-800 text-zinc-300 px-1.5 py-0.5 rounded font-mono">Deepgram</span>
          </div>
          <div className="text-base font-semibold text-zinc-100 font-mono">
            {formattedSttSeconds}s
          </div>
          <div className="text-[11px] text-zinc-500">
            Audio vAD continu
          </div>
        </div>

        {/* Column 3: TTS Text-to-Speech */}
        <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-3.5 space-y-1.5">
          <div className="flex justify-between items-center text-zinc-400 font-medium">
            <span>TTS Synthese</span>
            <span className="text-[10px] bg-zinc-800 text-zinc-300 px-1.5 py-0.5 rounded font-mono">ElevenLabs</span>
          </div>
          <div className="text-base font-semibold text-zinc-100 font-mono">
            {telemetry.ttsCharacters.toLocaleString()} chars
          </div>
          <div className="text-[11px] text-zinc-500">
            WebSocket streaming PCM
          </div>
        </div>

        {/* Column 4: SIP Telephony */}
        <div className="bg-zinc-950/80 border border-zinc-800/80 rounded-lg p-3.5 space-y-1.5">
          <div className="flex justify-between items-center text-zinc-400 font-medium">
            <span>SIP / Carrier</span>
            <span className="text-[10px] bg-zinc-800 text-zinc-300 px-1.5 py-0.5 rounded font-mono">Twilio / Telnyx</span>
          </div>
          <div className="text-base font-semibold text-zinc-100 font-mono">
            {formattedSipSeconds}s
          </div>
          <div className="text-[11px] text-emerald-400 flex items-center gap-1 font-mono">
            <span>[Active]</span> {telemetry.activeCallsCount} appels actifs
          </div>
        </div>
      </div>
    </div>
  );
};
