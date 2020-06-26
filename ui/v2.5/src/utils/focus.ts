import { useRef } from "react"

const useFocus = () => {
    const htmlElRef = useRef<any>();
    const setFocus = () => {
        const currentEl = htmlElRef.current
		currentEl && currentEl.focus();
    }

    return [ htmlElRef, setFocus ] as const;
}

export default useFocus;