const entries = []
{{range $i,$v := . }}
import entry{{$i}} from "flow:./{{$v.Name.Get}}"
entries.push({main:entry{{$i}}, name:"{{$v.Name.String}}" })
{{end}}
export default function (nodes) {
    const walk = (step, nodes) => {
        if (step >= entries.length) {
            return
        }
        const entry = entries[step]
        if (!entry?.main?.call) {
            console.warn(`flow: entry ${entry.name} is undefined, skip to the next`)
            walk(step + 1, nodes)
            return
        }
        this.FLOW_PIPE_NAME = entry.name
        entry.main.call(this, nodes, next(step + 1))
    }
    const next = (n = 0) => {
        const f = walk.bind(this, n)
        f.skip = next.bind(this, n + 1)
        return f
    }
    if (entries.length == 0) {
        throw new Error("flow: entries is not defined")
    }
    this.FLOW_PIPE_ORIGIN_UUID = entries[0].uuid;
    this.FLOW_PIPE_ORIGIN_NAME = entries[0].name;
    walk(0, nodes)
}