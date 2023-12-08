
#ifdef GL_ES

precision highp float;
#define IN varying
#define OUT out
#define TEXTURE texture2D

#else

#define IN in
#define OUT out
#define TEXTURE texture

#endif

uniform float time;
uniform vec2 resolution;


float rand(float n){return fract(sin(n) * 43758.5453123);}

float noise(float p){
    float fl = floor(p);
    float fc = fract(p);
    return mix(rand(fl), rand(fl + 1.0), fc);
}


vec3 rgb2hsb( in vec3 c ){
    vec4 K = vec4(0.0, -1.0 / 3.0, 2.0 / 3.0, -1.0);
    vec4 p = mix(vec4(c.bg, K.wz),
    vec4(c.gb, K.xy),
    step(c.b, c.g));
    vec4 q = mix(vec4(p.xyw, c.r),
    vec4(c.r, p.yzx),
    step(p.x, c.r));
    float d = q.x - min(q.w, q.y);
    float e = 1.0e-10;
    return vec3(abs(q.z + (q.w - q.y) / (6.0 * d + e)),
    d / (q.x + e),
    q.x);
}


vec3 hsb2rgb( in vec3 c ){
    vec3 rgb = clamp(abs(mod(c.x*6.0+vec3(0.0,4.0,2.0),
    6.0)-3.0)-1.0,
    0.0,
    1.0 );
    rgb = rgb*rgb*(3.0-2.0*rgb);
    return c.z * mix(vec3(1.0), rgb, c.y);
}


void main(){

    vec2 uv = gl_FragCoord.xy / resolution.xy;

    float t = time / 5.0;

    vec3 color = vec3(0.0);

    float uv_col = floor(uv.x * 14.0) / 14.0;

    float iterations = floor(0.5 * 8.0) + 1.0;

    for(float i = 0.0; i < iterations; i++) {
        float position = (
            noise(i + t) * 0.2 +
            sin(i * 2.1 + t) * 0.2 +
            sin(i * 1.72 + t*1.121) * 0.2 +
            sin(i * 2.221 + t*0.437) * 0.2 +
            sin(i * 3.1122 + t*4.269) * 0.2
        ) + 0.5;
        float dist = abs(uv_col - position);
        float pct = 1.0 - smoothstep(0.0, 1.0/11.0, dist);
        float c = (
            noise(i + t) * 0.2 +
            sin(i * 3.1 + t) * 0.2 +
            sin(i * 5.21 + t*1.69) * 0.2 +
            sin(i * 1.0845 + t*0.278) * 0.2 +
            sin(i * 2.885 + t*4.201) * 0.2
        ) + 0.5;
        color += hsb2rgb(vec3(c,1.0,1.0 * pct)) * 0.8;
    }

    gl_FragColor = vec4(color, 1.0);


}
