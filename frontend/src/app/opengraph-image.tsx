import { ImageResponse } from "next/og";

export const alt = "UBC Esports Memberships";
export const size = {
  width: 1200,
  height: 630,
};
export const contentType = "image/png";

export default function OpenGraphImage() {
  return new ImageResponse(
    <div
      style={{
        alignItems: "center",
        background: "#070b14",
        color: "#f4f7ff",
        display: "flex",
        height: "100%",
        justifyContent: "center",
        padding: "64px",
        position: "relative",
        width: "100%",
      }}
    >
      <div
        style={{
          background: "#1f64f0",
          height: "12px",
          left: 0,
          position: "absolute",
          top: 0,
          width: "100%",
        }}
      />
      <div
        style={{
          border: "2px solid #25314a",
          display: "flex",
          flexDirection: "column",
          height: "100%",
          justifyContent: "space-between",
          padding: "58px 64px",
          width: "100%",
        }}
      >
        <div style={{ color: "#73a2ff", display: "flex", fontSize: 28, letterSpacing: 8 }}>
          UBCEA
        </div>
        <div style={{ display: "flex", flexDirection: "column" }}>
          <div style={{ display: "flex", fontSize: 72, fontWeight: 700, lineHeight: 1.05 }}>
            UBC Esports Memberships
          </div>
          <div style={{ color: "#b8c2d9", display: "flex", fontSize: 31, marginTop: 24 }}>
            Compare passes, pricing, and member benefits.
          </div>
        </div>
        <div style={{ color: "#7f8ba5", display: "flex", fontSize: 24 }}>app.ubcesports.ca</div>
      </div>
    </div>,
    size,
  );
}
