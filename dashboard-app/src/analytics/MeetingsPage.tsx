import React, { useState } from 'react';

export const MeetingsPage: React.FC = () => {
  const [videoResolution, setVideoResolution] = useState<'720p' | '1080p'>('1080p');
  const [botMode, setBotMode] = useState<'rtms_native' | 'signed_in_bot'>('signed_in_bot');
  const [multiTrackIsolation, setMultiTrackIsolation] = useState<boolean>(true);

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="border-b border-slate-800 pb-4 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Suivi des Reunions Visio & Moteur Bot</h1>
          <p className="text-sm text-slate-400">Bots Zoom (RTMS), Google Meet & Microsoft Teams: temps de parole, diarisation et capture audio/video isolee.</p>
        </div>
        <span className="text-xs font-mono bg-cyan-950 text-cyan-400 border border-cyan-800 px-3 py-1 rounded-full font-semibold self-start md:self-auto">
          Moteur Visio Native v2.4
        </span>
      </div>

      {/* Configuration & Capabilities Banner */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
        <h3 className="font-bold text-white text-base">Configuration des Bots Visio Multi-Plateformes</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
          <div className="bg-slate-950 p-4 rounded-lg border border-slate-800 space-y-2">
            <label className="font-semibold text-slate-300 block">Mode de Connexion Bot</label>
            <select
              value={botMode}
              onChange={(e) => setBotMode(e.target.value as 'rtms_native' | 'signed_in_bot')}
              className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-slate-200 focus:outline-none focus:border-cyan-500"
            >
              <option value="signed_in_bot">Bot Identifie (Zoom / Teams / Meet)</option>
              <option value="rtms_native">Flux Direct RTMS / Native Media API</option>
            </select>
            <p className="text-[11px] text-slate-500">Capture le flux audio et video directement avec identite personnalisee du tenant.</p>
          </div>

          <div className="bg-slate-950 p-4 rounded-lg border border-slate-800 space-y-2">
            <label className="font-semibold text-slate-300 block">Resolution Video Render</label>
            <select
              value={videoResolution}
              onChange={(e) => setVideoResolution(e.target.value as '720p' | '1080p')}
              className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-slate-200 focus:outline-none focus:border-cyan-500"
            >
              <option value="720p">720p HD (Economique)</option>
              <option value="1080p">1080p Full HD (Haute Fidelite)</option>
            </select>
            <p className="text-[11px] text-slate-500">Qualite de rendu video pour la retransmission et l'analyse visuelle.</p>
          </div>

          <div className="bg-slate-950 p-4 rounded-lg border border-slate-800 space-y-2">
            <label className="font-semibold text-slate-300 block">Isolation Audio Multi-Pistes</label>
            <div className="flex items-center justify-between pt-1">
              <span className="text-slate-400">Pistes Pcm Separees</span>
              <input
                type="checkbox"
                checked={multiTrackIsolation}
                onChange={(e) => setMultiTrackIsolation(e.target.checked)}
                className="h-4 w-4 rounded bg-slate-900 border-slate-700 text-cyan-500 focus:ring-cyan-500"
              />
            </div>
            <p className="text-[11px] text-slate-500">Isole le flux micro de chaque participant pour une diarisation parfaite.</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Talk Time & Diarization */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Temps de Parole & Diarisation</h3>
          <p className="text-xs text-slate-400">Repartition du temps d'expression durant les reunions visio.</p>

          <div className="space-y-3 text-xs">
            <div>
              <div className="flex justify-between font-medium mb-1">
                <span className="text-slate-300">Interlocuteurs (Prospects)</span>
                <span className="text-cyan-400 font-bold">62%</span>
              </div>
              <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                <div className="bg-cyan-400 h-full" style={{ width: '62%' }}></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between font-medium mb-1">
                <span className="text-slate-300">Agent IA Patter Engine</span>
                <span className="text-indigo-400 font-bold">38%</span>
              </div>
              <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                <div className="bg-indigo-400 h-full" style={{ width: '38%' }}></div>
              </div>
            </div>
          </div>
        </div>

        {/* Objections & Copilot Suggestions */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Objections & Suggestions Copilote</h3>
          <p className="text-xs text-slate-400">Liste des points complexes souleves et reponses guidees par l'IA.</p>

          <div className="space-y-3 text-xs">
            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1">
              <div className="text-amber-300 font-semibold">"Pouvons-nous herberger les donnees On-Premise ?"</div>
              <div className="text-slate-400">Suggestion Copilote: "Proposer le deploiement hybride sur VPC dedie AWS/GCP."</div>
            </div>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 space-y-1">
              <div className="text-amber-300 font-semibold">"Quel est le delai de mise en production ?"</div>
              <div className="text-slate-400">Suggestion Copilote: "Indiquer un deploiement standard en 48h avec connecteur CRM."</div>
            </div>
          </div>
        </div>

        {/* Post-Meeting Action Items */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <h3 className="font-bold text-white text-base">Suivi des Actions Post-Reunion</h3>
          <p className="text-xs text-slate-400">Taux d'execution des compte-rendus et taches attribuees.</p>

          <div className="space-y-3 text-xs">
            <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800">
              <span className="text-slate-300">Compte-Rendus Envoyes Automatiquement</span>
              <span className="text-emerald-400 font-bold bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">98.5%</span>
            </div>

            <div className="flex justify-between items-center bg-slate-950 p-3 rounded-lg border border-slate-800">
              <span className="text-slate-300">Taches CRM / Follow-up creees</span>
              <span className="text-cyan-400 font-bold bg-cyan-500/10 px-2 py-0.5 rounded border border-cyan-500/20">142 taches</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
